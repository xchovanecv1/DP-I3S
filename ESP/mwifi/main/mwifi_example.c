/*
 * ESPRESSIF MIT License
 *
 * Copyright (c) 2018 <ESPRESSIF SYSTEMS (SHANGHAI) PTE LTD>
 *
 * Permission is hereby granted for use on all ESPRESSIF SYSTEMS products, in which case,
 * it is free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished
 * to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/gpio.h"

#include "esp_system.h"
#include "esp_log.h"
#include "driver/uart.h"
#include "string.h"
#include "mdf_common.h"
#include "mwifi.h"

#include "cJSON.h"

// #define MEMORY_DEBUG

static int g_sockfd = -1;
static const char *TAG = "i3s_gateway";

#define TXD_PIN (GPIO_NUM_17)
#define RXD_PIN (GPIO_NUM_16)

static const int RX_BUF_SIZE = 1024;

#define GPIO_INPUT_IO_2 24 // Data 1
#define GPIO_INPUT_IO_3 13
#define GPIO_INPUT_PIN_SEL2 ((1ULL << GPIO_INPUT_IO_2) | (1ULL << GPIO_INPUT_IO_3))
#define GPIO_INPUT_IO_0 4
#define GPIO_INPUT_IO_1 5
#define GPIO_INPUT_PIN_SEL ((1ULL << GPIO_INPUT_IO_0) | (1ULL << GPIO_INPUT_IO_1))
#define ESP_INTR_FLAG_DEFAULT 0

typedef struct
{
    //xQueueHandle gpio_evt_queue;
    volatile unsigned long _cardTempHigh;
    volatile unsigned long long _cardTemp;
    volatile unsigned int _lastWiegand;
    volatile uint8_t _cardRead;
    unsigned long long _code;
    volatile int _bitCount;
    int _wiegandType;

} __wIF;

__wIF wiegandTerm[2];

xQueueHandle RFID_ACCESS_Q;
xQueueHandle TAG_ACCESS_Q;

#define HEARTBEAT_TIME 15000
unsigned int lastHearhtb = 0;
bool hbPending = false;

#define MIN_NEW_CARD_INTERVAL 3000
unsigned int lastTagTime = 0;
unsigned long long lastTagId = 0;

unsigned long long messageID = 1;

#define MAX_REPOST_DEPTH 20
#define MAX_REPOST_TIME 30000

unsigned long long cardIds[MAX_REPOST_DEPTH];
unsigned long long cardMessageIds[MAX_REPOST_DEPTH];
unsigned int cardRepostTime[MAX_REPOST_DEPTH];

void addPendingId(unsigned long long mid, unsigned long long cid)
{
    for (int i = 0; i < MAX_REPOST_DEPTH; i++)
    {
        if (cardMessageIds[i] == 0)
        {
            cardIds[i] = cid;
            cardMessageIds[i] = mid;
            cardRepostTime[i] = xTaskGetTickCount() + MAX_REPOST_TIME;
            break;
        }
    }
}

void confirmPendingID(unsigned long long id)
{
    unsigned int now = xTaskGetTickCount();
    for (int i = 0; i < MAX_REPOST_DEPTH; i++)
    {
        if (cardMessageIds[i] == id)
        {
            cardMessageIds[i] = 0;
            cardRepostTime[i] = 0;
            cardIds[i] = 0;
            break;
        }
    }
}

void checkPendingID()
{
    unsigned int now;
    for (;;)
    {
        now = xTaskGetTickCount();
        for (int i = 0; i < MAX_REPOST_DEPTH; i++)
        {
            if (cardMessageIds[i] != 0 && now >= cardRepostTime[i])
            {
                ESP_LOGI("TRANSIT REPOST", "Card transit %llu failed to response, resending..", cardMessageIds[i]);

                xQueueSend(TAG_ACCESS_Q, (void *)&(cardIds[i]), (TickType_t)0);
                cardMessageIds[i] = 0;
                cardRepostTime[i] = 0;
                cardIds[i] = 0;
            }
        }
        vTaskDelay(128 / portTICK_RATE_MS);
    }
}

void uart_init(void)
{
    const uart_config_t uart_config = {
        .baud_rate = 9600,
        .data_bits = UART_DATA_8_BITS,
        .parity = UART_PARITY_DISABLE,
        .stop_bits = UART_STOP_BITS_1,
        .flow_ctrl = UART_HW_FLOWCTRL_DISABLE};
    uart_param_config(UART_NUM_1, &uart_config);
    uart_set_pin(UART_NUM_1, TXD_PIN, RXD_PIN, UART_PIN_NO_CHANGE, UART_PIN_NO_CHANGE);
    // We won't use a buffer for sending data.
    uart_driver_install(UART_NUM_1, RX_BUF_SIZE * 2, 0, 0, NULL, 0);
}

/*
    Handling for dual card uart reader
 */

typedef enum dRFID_CardType
{
    UNKNOWN,
    RFID,
    EM4100,
} RFID_CARD_TYPES;

void dRFID_parseOutput(uint8_t *data, int len)
{
    if (len > 8)
    {
        unsigned long long cardId = 0;
        static const char *dRFID_parseOutput_TAG = "dRFID_parseOutput";
        // First data byte is not message start
        for (uint8_t i = 0; i < len; i++)
        {
            ESP_LOGI(dRFID_parseOutput_TAG, "%x ", data[i]);
        }

        if (data[0] != 0x02)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Card data not starting with 0x02 instead %02x", data[0]);
            return;
        }
        //Check message length
        if (data[1] != len)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Card data not correct length of %d instead %d", data[1], len);
            return;
        }

        //Check ending byte
        if (data[len - 1] != 0x03)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Card data not ending with 0x03 instead %02x", data[len - 1]);
            return;
        }
        RFID_CARD_TYPES cT = data[2];
        //Determine check logic based on card type
        switch (cT)
        {
        case EM4100:
            ESP_LOGI(dRFID_parseOutput_TAG, "125 Data start");
            for (uint8_t i = 3; i < 8; i++)
            {
                cardId |= (data[i] & 0xFF);
                cardId <<= 8;
                ESP_LOGI(dRFID_parseOutput_TAG, "%x ", data[i]);
            }
            cardId >>= 8;

            if (cardId != lastTagId || xTaskGetTickCount() - lastTagTime > MIN_NEW_CARD_INTERVAL)
            {
                lastTagId = cardId;
                lastTagTime = xTaskGetTickCount();
                xQueueSend(TAG_ACCESS_Q, (void *)&cardId, (TickType_t)0);
                ESP_LOGI(dRFID_parseOutput_TAG, "time : %d", lastTagTime);
                ESP_LOGI(dRFID_parseOutput_TAG, "Card id: %llu", cardId);
            }
            else if (cardId == lastTagId)
            {
                lastTagTime = xTaskGetTickCount();
            }
            break;

        case RFID:
            ESP_LOGI(dRFID_parseOutput_TAG, "13.56 Data start");
            for (uint8_t i = 3; i < 7; i++)
            {
                cardId |= (data[i] & 0xFF);
                cardId <<= 8;
                ESP_LOGI(dRFID_parseOutput_TAG, "%x ", data[i]);
            }
            cardId >>= 8;
            xQueueSend(RFID_ACCESS_Q, (void *)&cardId, (TickType_t)0);
            ESP_LOGI(dRFID_parseOutput_TAG, "Card id: %llu", cardId);
            break;

        default: // UNKNOWN, just dump data
            for (uint8_t i = 0; i < len; i++)
            {
                ESP_LOGI(dRFID_parseOutput_TAG, "%x ", data[i]);
            }
            break;
        }
    }
}

void RFID125_parseOutput(uint8_t *data, int len)
{
    unsigned long long cardId = 0;
    uint32_t tag_id;
    uint8_t checksum;
    char buff[14];
    if (len > 13)
    {
        if (len > 14)
        {
            len = 14;
        }

        static const char *dRFID_parseOutput_TAG = "dRFID_parseOutput";
        // First data byte is not message start
        for (uint8_t i = 0; i < len; i++)
        {
            ESP_LOGI(dRFID_parseOutput_TAG, "%x ", data[i]);
            buff[i] = data[i];
        }

        if (data[0] != 0x02)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Card data not starting with 0x02 instead %02x", data[0]);
            return;
        }

        //Check ending byte
        if (data[len - 1] != 0x03)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Card data not ending with 0x03 instead %02x", data[len - 1]);
            return;
        }
        RFID_CARD_TYPES cT = 2;

        buff[13] = 0;
        checksum = strtol(buff + 11, NULL, 16);
        /* add null and parse tag_id */
        buff[11] = 0;
        cardId = strtol(buff + 3, NULL, 16);
        /* add null and parse version (needs to be xored with checksum) */
        buff[3] = 0;
        checksum ^= strtol(buff + 1, NULL, 16);
        //Determine check logic based on card type

        cardId += 433791696896;

        for (uint8_t i = 0; i < 32; i += 8)
        {
            checksum ^= ((cardId >> i) & 0xFF);
        }
        if (checksum)
        {
            ESP_LOGW(dRFID_parseOutput_TAG, "Invalid checksum %d", checksum);
        }

        if (cardId != lastTagId || xTaskGetTickCount() - lastTagTime > MIN_NEW_CARD_INTERVAL)
        {
            lastTagId = cardId;
            lastTagTime = xTaskGetTickCount();
            xQueueSend(TAG_ACCESS_Q, (void *)&cardId, (TickType_t)0);
            ESP_LOGI(dRFID_parseOutput_TAG, "time : %d", lastTagTime);
            ESP_LOGI(dRFID_parseOutput_TAG, "Card id: %llu", cardId);
        }
        else if (cardId == lastTagId)
        {
            lastTagTime = xTaskGetTickCount();
        }

        ESP_LOGI(dRFID_parseOutput_TAG, "Card id: %llu", cardId);
    }
}

static void uart_rx_task(void *arg)
{
    static const char *RX_TASK_TAG = "RX_TASK";
    esp_log_level_set(RX_TASK_TAG, ESP_LOG_INFO);
    uint8_t *data = (uint8_t *)malloc(RX_BUF_SIZE + 1);
    while (1)
    {
        const int rxBytes = uart_read_bytes(UART_NUM_1, data, RX_BUF_SIZE, 1000 / portTICK_RATE_MS);
        if (rxBytes > 0)
        {
            dRFID_parseOutput(data, rxBytes);
            RFID125_parseOutput(data, rxBytes);
        }
    }
    free(data);
}
/*

/////////////////////////////////////////////////////////////////////////////////////////

 */
static void setup_Weigand(__wIF *wif, uint8_t D0, uint8_t D1);

static void IRAM_ATTR gpio_isr_handler(void *arg)
{
    uint32_t gpio_num = (uint32_t)arg;
    __wIF *wif;
    if (gpio_num == GPIO_INPUT_IO_2 || gpio_num == GPIO_INPUT_IO_3)
    {
        wif = &(wiegandTerm[0]);
    }
    else
    {
        wif = &(wiegandTerm[1]);
    }

    if (wif->_bitCount > 34)
    { // something fishy - return
        return;
    }
    //ets_printf("GPIO bef %llu %llu %d\n", wif->_cardTempHigh, wif->_cardTemp, wif->_bitCount);
    // ets_printf("Bitc %d\n", wif->_bitCount);
    /*if (wif->_bitCount>31)			// If bit count more than 31, process high bits
	{
		wif->_cardTempHigh |= ((0x80000000 & wif->_cardTemp)>>31);	//	shift value to high bits
		wif->_cardTempHigh <<= 1;
		if(gpio_num == GPIO_INPUT_IO_0 || gpio_num == GPIO_INPUT_IO_2) { //D1
			wif->_cardTemp |= 1;
            ets_printf("1\n");
		} else {

            ets_printf("0\n");
        }
		wif->_cardTemp <<=1;
	}
	else
	{*/
    if (gpio_num == GPIO_INPUT_IO_0 || gpio_num == GPIO_INPUT_IO_2)
    { //D1
        wif->_cardTemp |= 1;
        ets_printf("1");
    }
    else
    {

        ets_printf("0");
    }
    wif->_cardTemp <<= 1; // D0 represent binary 0, so just left shift card data
    //}
    wif->_bitCount++; // Increament bit count for Interrupt connected to D0
    wif->_lastWiegand = xTaskGetTickCountFromISR();
    //printf("Int %zu", gpio_num);
    //xQueueSendFromISR(gpio_evt_queue, &gpio_num, NULL);
}

static void wiegand_check_trans(void *arg)
{

    __wIF *wif1 = &(wiegandTerm[0]);
    __wIF *wif2 = &(wiegandTerm[1]);
    uint8_t dir = 0;
    for (;;)
    {
        if (wif1->_code != 0 && wif2->_code != 0)
        {
            if (wif1->_lastWiegand > wif2->_lastWiegand)
            {
                dir = 0;
            }
            else
            {
                dir = 1;
            }

            if (wif1->_code == wif2->_code)
            {
                printf("Prechod %d jednotky\n", dir);
            }
            wif1->_lastWiegand = 0;
            wif1->_cardTempHigh = 0;
            wif1->_cardTemp = 0;
            wif1->_cardRead = 0;
            wif1->_code = 0;
            wif1->_wiegandType = 0;
            wif1->_bitCount = 0;

            wif2->_lastWiegand = 0;
            wif2->_cardTempHigh = 0;
            wif2->_cardTemp = 0;
            wif2->_cardRead = 0;
            wif2->_code = 0;
            wif2->_wiegandType = 0;
            wif2->_bitCount = 0;
            /*
			wif1->_code = 0;
			wif2->_code = 0;*/
        }
        vTaskDelay(128 / portTICK_RATE_MS);
    }
}

static void gateway_heartbeat(void *arg)
{
    unsigned int tm = 0;
    for (;;)
    {
        if (!mwifi_is_connected())
        {
            vTaskDelay(500 / portTICK_RATE_MS);
            continue;
        }

        tm = xTaskGetTickCount();
        if (!hbPending && (tm - lastHearhtb) > HEARTBEAT_TIME)
        {
            lastHearhtb = tm;
            hbPending = true;
        }
        vTaskDelay(2000 / portTICK_RATE_MS);
    }
}

static void gpio_task_example(void *arg)
{
    __wIF *wif = (__wIF *)arg;
    uint8_t p_even, p_odd;
    uint8_t p_even_c = 0;
    uint8_t p_odd_c = 0;
    //unsigned long cardID = 0;
    for (;;)
    {

        if ((xTaskGetTickCount() - wif->_lastWiegand) > 5) // if no more signal coming through after 25ms
        {
            //printf("Going to chceck weigand code %d\n", wif->_bitCount==26);
            //printf("AFTER %llu %d\n", wif->_cardTemp, wif->_bitCount);
            if (wif->_bitCount == 26)
            { // EM tag
                wif->_code = (wif->_cardTemp & 0x1FFFFFE) >> 1;
                wif->_wiegandType = 26;
                wif->_cardRead = 1;
            }
            if (wif->_bitCount == 34) // Mifare
            {
                wif->_wiegandType = 34;
                //wif->_cardTempHigh = wif->_cardTempHigh & 0x03;				// only need the 2 LSB of the codehigh

                //printf("CTH %lu \n", wif->_cardTempHigh);
                //wif->_cardTempHigh <<= 29;							// shift 2 LSB to MSB
                //printf("CTH %lu \n", wif->_cardTempHigh);
                wif->_cardTemp >>= 1;
                p_even_c = 0;
                p_odd_c = 0;

                p_even = wif->_cardTemp & 0x1;
                p_odd = (wif->_cardTemp & 0x100000000) >> 33;

                wif->_cardTemp >>= 1;

                printf("DTA %llu e:%d o:%d \n", wif->_cardTemp, p_even, p_odd);
                wif->_code = wif->_cardTemp & 0xFFFFFFFF; // take only first 32 bites

                for (uint8_t i = 0; i < 32; i++)
                {
                    if (i < 16)
                    { // even
                        p_even_c += (uint8_t)(wif->_cardTemp & 1);
                    }
                    else
                    {
                        p_odd_c += (uint8_t)(wif->_cardTemp & 1);
                    }
                    wif->_cardTemp >>= 1;
                }
                if ((~(p_even_c % 2) & 1) == p_even && (p_odd_c % 2) == p_odd)
                {
                    wif->_cardRead = 1;
                    printf("CDE %llu\n", wif->_code);
                }
                printf("ec %d %d\n", p_even_c, (~(p_even_c % 2) & 1));
                printf("oc %d %d\n", p_odd_c, (p_odd_c % 2));
            }

            if (wif->_cardRead)
            {
                if (wif->_wiegandType == 34 && RFID_ACCESS_Q != 0)
                {
                    wif->_cardRead = 0;
                    xQueueSend(RFID_ACCESS_Q, (void *)&(wif->_code), (TickType_t)0);
                }
                else if (wif->_wiegandType == 26 && TAG_ACCESS_Q != 0)
                {
                    wif->_cardRead = 0;
                    if (wif->_code != lastTagId || xTaskGetTickCount() - lastTagTime > MIN_NEW_CARD_INTERVAL)
                    {
                        lastTagId = wif->_code;
                        lastTagTime = xTaskGetTickCount();
                        xQueueSend(TAG_ACCESS_Q, (void *)&(wif->_code), (TickType_t)0);
                    }
                }
            }

            if (wif->_bitCount >= 34)
            {
                wif->_lastWiegand = 0;
                wif->_cardTempHigh = 0;
                wif->_cardTemp = 0;
                wif->_cardRead = 0;
                wif->_wiegandType = 0;
                wif->_bitCount = 0;
                printf("Wrong data, restart structure\n");
                wif->_code = 1569850504;
                xQueueSend(RFID_ACCESS_Q, (void *)&(wif->_code), (TickType_t)0);
            }
        }
        vTaskDelay(100 / portTICK_RATE_MS);
    }
}

void tcp_client_write_task(void *arg)
{
    mdf_err_t ret = MDF_OK;
    char *data = MDF_CALLOC(1, MWIFI_PAYLOAD_LEN);
    size_t size = MWIFI_PAYLOAD_LEN;
    uint8_t src_addr[MWIFI_ADDR_LEN] = {0x0};
    mwifi_data_type_t data_type = {0x0};

    MDF_LOGI("TCP client write task is running");

    while (mwifi_is_connected())
    {
        if (g_sockfd == -1)
        {
            vTaskDelay(500 / portTICK_RATE_MS);
            continue;
        }

        size = MWIFI_PAYLOAD_LEN - 1;
        memset(data, 0, MWIFI_PAYLOAD_LEN);
        ret = mwifi_root_read(src_addr, &data_type, data, &size, portMAX_DELAY);
        MDF_ERROR_CONTINUE(ret != MDF_OK, "<%s> mwifi_root_read", mdf_err_to_name(ret));

        char *json_data = NULL;
        int json_size = asprintf(&json_data, "{\"addr\":\"" MACSTR "\",\"data\":%s}",
                                 MAC2STR(src_addr), data);

        MDF_LOGI("TCP write, size: %d, data: %s", json_size, json_data);
        ret = write(g_sockfd, json_data, json_size);
        MDF_FREE(json_data);
        MDF_ERROR_CONTINUE(ret <= 0, "<%s> TCP write", strerror(errno));
    }

    MDF_LOGI("TCP client write task is exit");

    close(g_sockfd);
    MDF_FREE(data);
    vTaskDelete(NULL);
}

/**
 * @brief Create a tcp client
 */
static int socket_tcp_client_create(const char *ip, uint16_t port)
{
    MDF_PARAM_CHECK(ip);

    MDF_LOGI("Create a tcp client, ip: %s, port: %d", ip, port);

    mdf_err_t ret = ESP_OK;
    int sockfd = -1;
    struct sockaddr_in server_addr = {
        .sin_family = AF_INET,
        .sin_port = htons(port),
        .sin_addr.s_addr = inet_addr(ip),
    };

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    MDF_ERROR_GOTO(sockfd < 0, ERR_EXIT, "socket create, sockfd: %d", sockfd);

    ret = connect(sockfd, (struct sockaddr *)&server_addr, sizeof(struct sockaddr_in));
    MDF_ERROR_GOTO(ret < 0, ERR_EXIT, "socket connect, ret: %d, ip: %s, port: %d",
                   ret, ip, port);
    return sockfd;

ERR_EXIT:

    if (sockfd != -1)
    {
        close(sockfd);
    }

    return -1;
}

void tcp_client_read_task(void *arg)
{
    mdf_err_t ret = MDF_OK;
    char *data = MDF_MALLOC(MWIFI_PAYLOAD_LEN);
    size_t size = MWIFI_PAYLOAD_LEN;
    uint8_t dest_addr[MWIFI_ADDR_LEN] = {0x0};
    mwifi_data_type_t data_type = {0x0};

    MDF_LOGI("TCP client read task is running");

    while (mwifi_is_connected())
    {
        if (g_sockfd == -1)
        {
            g_sockfd = socket_tcp_client_create(CONFIG_SERVER_IP, CONFIG_SERVER_PORT);

            if (g_sockfd == -1)
            {
                vTaskDelay(500 / portTICK_RATE_MS);
                continue;
            }
        }

        memset(data, 0, MWIFI_PAYLOAD_LEN);
        ret = read(g_sockfd, data, size);
        MDF_LOGI("TCP read, %d, size: %d, data: %s", g_sockfd, size, data);

        if (ret <= 0)
        {
            MDF_LOGW("<%s> TCP read", strerror(errno));
            close(g_sockfd);
            g_sockfd = -1;
            continue;
        }

        cJSON *pJson = NULL;
        cJSON *pSub = NULL;

        pJson = cJSON_Parse(data);
        MDF_ERROR_CONTINUE(!pJson, "cJSON_Parse, data format error");

        pSub = cJSON_GetObjectItem(pJson, "addr");

        if (!pSub)
        {
            MDF_LOGW("cJSON_GetObjectItem, Destination address not set");
            cJSON_Delete(pJson);
            continue;
        }

        /**
         * @brief  Convert mac from string format to binary
         */
        do
        {
            uint32_t mac_data[MWIFI_ADDR_LEN] = {0};
            sscanf(pSub->valuestring, MACSTR,
                   mac_data, mac_data + 1, mac_data + 2,
                   mac_data + 3, mac_data + 4, mac_data + 5);

            for (int i = 0; i < MWIFI_ADDR_LEN; i++)
            {
                dest_addr[i] = mac_data[i];
            }
        } while (0);

        pSub = cJSON_GetObjectItem(pJson, "data");

        if (!pSub)
        {
            MDF_LOGW("cJSON_GetObjectItem, Failed to get data");
            cJSON_Delete(pJson);
            continue;
        }

        char *json_data = cJSON_PrintUnformatted(pSub);

        ret = mwifi_root_write(dest_addr, 1, &data_type, json_data, strlen(json_data), true);
        MDF_ERROR_CONTINUE(ret != MDF_OK, "<%s> mwifi_root_write", mdf_err_to_name(ret));

        MDF_FREE(json_data);
        cJSON_Delete(pJson);
    }

    MDF_LOGI("TCP client read task is exit");

    close(g_sockfd);
    g_sockfd = -1;
    MDF_FREE(data);
    vTaskDelete(NULL);
}

static void node_read_task(void *arg)
{
    mdf_err_t ret = MDF_OK;
    cJSON *pJson = NULL;
    cJSON *pSub = NULL;
    char *data = MDF_MALLOC(MWIFI_PAYLOAD_LEN);
    size_t size = MWIFI_PAYLOAD_LEN;
    mwifi_data_type_t data_type = {0x0};
    uint8_t src_addr[MWIFI_ADDR_LEN] = {0x0};

    MDF_LOGI("Note read task is running");

    for (;;)
    {
        if (!mwifi_is_connected())
        {
            vTaskDelay(500 / portTICK_RATE_MS);
            continue;
        }

        size = MWIFI_PAYLOAD_LEN;
        memset(data, 0, MWIFI_PAYLOAD_LEN);
        ret = mwifi_read(src_addr, &data_type, data, &size, portMAX_DELAY);
        MDF_ERROR_CONTINUE(ret != MDF_OK, "<%s> mwifi_read", mdf_err_to_name(ret));
        MDF_LOGD("Node receive: " MACSTR ", size: %d, data: %s", MAC2STR(src_addr), size, data);

        pJson = cJSON_Parse(data);
        MDF_ERROR_CONTINUE(!pJson, "cJSON_Parse, data format error, data: %s", data);

        pSub = cJSON_GetObjectItem(pJson, "mid");
        if (!pSub)
        {
            cJSON_Delete(pJson);
            continue;
        } else {
            confirmPendingID((unsigned long long)pSub->valueint);
        }

        pSub = cJSON_GetObjectItem(pJson, "auth");
        //MDF_LOGD("User authde: %d\n", pJson->string == NULL);
        //MDF_LOGD("User ret: %d\n", pSub);


        if (!pSub)
        {

            const char *error_ptr = cJSON_GetErrorPtr();
            if (error_ptr != NULL)
            {
                MDF_LOGD("Error before: %s\n", error_ptr);
            }
            cJSON_Delete(error_ptr);
            MDF_LOGW("cJSON_GetObjectItem, no card");
            cJSON_Delete(pJson);
            continue;
        }

        MDF_LOGD("User authd: %d", pSub->valueint);
        // Set otuput for buzzer
        gpio_set_level(CONFIG_LED_GPIO_NUM, pSub->valueint);

        cJSON_Delete(pJson);
    }

    MDF_LOGW("Note read task is exit");

    MDF_FREE(data);
    vTaskDelete(NULL);
}

static void node_write_task(void *arg)
{
    mdf_err_t ret = MDF_OK;
    int count = 0;
    size_t size = 0;
    char *data = NULL;
    mwifi_data_type_t data_type = {0x0};

    MDF_LOGI("NODE task is running");

    for (;;)
    {
        if (!mwifi_is_connected())
        {
            vTaskDelay(500 / portTICK_RATE_MS);
            continue;
        }

        size = 0;

        unsigned long long pxRxedMessage;
        if (hbPending)
        {
            size = asprintf(&data, "{\"tick\":%u}", xTaskGetTickCount());
            hbPending = false;
        }
        else if (RFID_ACCESS_Q != 0 || TAG_ACCESS_Q != 0)
        {
            // Receive a message on the created queue.  Block for 10 ticks if a
            // message is not immediately available.

            if (xQueueReceive(TAG_ACCESS_Q, &(pxRxedMessage), (TickType_t)10))
            {
                size = asprintf(&data, "{\"type\":1,\"card\":\"%llu\", \"mid\":%llu}", pxRxedMessage, messageID);
                addPendingId(messageID, pxRxedMessage);
                messageID++;
            }
            else if (xQueueReceive(RFID_ACCESS_Q, &(pxRxedMessage), (TickType_t)10))
            {
                size = asprintf(&data, "{\"type\":0,\"card\":\"%llu\", \"mid\":%llu}", pxRxedMessage, messageID);
                messageID++;
            }
        }

        if (size != 0)
        {
            MDF_LOGD("Node send, size: %d, data: %s", size, data);
            ret = mwifi_write(NULL, &data_type, data, size, true);
            MDF_FREE(data);
            MDF_ERROR_CONTINUE(ret != MDF_OK, "<%s> mwifi_write", mdf_err_to_name(ret));
        }

        vTaskDelay(30 / portTICK_RATE_MS);
    }

    MDF_LOGW("NODE task is exit");

    vTaskDelete(NULL);
}

/**
 * @brief Timed printing system information
 */
static void print_system_info_timercb(void *timer)
{
    uint8_t primary = 0;
    wifi_second_chan_t second = 0;
    mesh_addr_t parent_bssid = {0};
    uint8_t sta_mac[MWIFI_ADDR_LEN] = {0};
    mesh_assoc_t mesh_assoc = {0x0};
    wifi_sta_list_t wifi_sta_list = {0x0};

    esp_wifi_get_mac(ESP_IF_WIFI_STA, sta_mac);
    esp_wifi_ap_get_sta_list(&wifi_sta_list);
    esp_wifi_get_channel(&primary, &second);
    esp_wifi_vnd_mesh_get(&mesh_assoc);
    esp_mesh_get_parent_bssid(&parent_bssid);

    MDF_LOGI("System information, channel: %d, layer: %d, self mac: " MACSTR ", parent bssid: " MACSTR
             ", parent rssi: %d, node num: %d, free heap: %u",
             primary,
             esp_mesh_get_layer(), MAC2STR(sta_mac), MAC2STR(parent_bssid.addr),
             mesh_assoc.rssi, esp_mesh_get_total_node_num(), esp_get_free_heap_size());

    for (int i = 0; i < wifi_sta_list.num; i++)
    {
        MDF_LOGI("Child mac: " MACSTR, MAC2STR(wifi_sta_list.sta[i].mac));
    }

#ifdef MEMORY_DEBUG
    if (!heap_caps_check_integrity_all(true))
    {
        MDF_LOGE("At least one heap is corrupt");
    }

    mdf_mem_print_heap();
    mdf_mem_print_record();
#endif /**< MEMORY_DEBUG */
}

static mdf_err_t wifi_init()
{
    mdf_err_t ret = nvs_flash_init();
    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();

    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
        MDF_ERROR_ASSERT(nvs_flash_erase());
        ret = nvs_flash_init();
    }

    MDF_ERROR_ASSERT(ret);

    tcpip_adapter_init();
    MDF_ERROR_ASSERT(esp_event_loop_init(NULL, NULL));
    MDF_ERROR_ASSERT(esp_wifi_init(&cfg));
    MDF_ERROR_ASSERT(esp_wifi_set_storage(WIFI_STORAGE_FLASH));
    MDF_ERROR_ASSERT(esp_wifi_set_mode(WIFI_MODE_STA));
    MDF_ERROR_ASSERT(esp_wifi_set_ps(WIFI_PS_NONE));
    MDF_ERROR_ASSERT(esp_mesh_set_6m_rate(false));
    MDF_ERROR_ASSERT(esp_wifi_start());

    return MDF_OK;
}

/**
 * @brief All module events will be sent to this task in esp-mdf
 *
 * @Note:
 *     1. Do not block or lengthy operations in the callback function.
 *     2. Do not consume a lot of memory in the callback function.
 *        The task memory of the callback function is only 4KB.
 */
static mdf_err_t event_loop_cb(mdf_event_loop_t event, void *ctx)
{
    MDF_LOGI("event_loop_cb, event: %d", event);

    switch (event)
    {
    case MDF_EVENT_MWIFI_STARTED:
        MDF_LOGI("MESH is started");
        break;

    case MDF_EVENT_MWIFI_PARENT_CONNECTED:
        MDF_LOGI("Parent is connected on station interface");
        break;

    case MDF_EVENT_MWIFI_PARENT_DISCONNECTED:
        MDF_LOGI("Parent is disconnected on station interface");
        break;

    case MDF_EVENT_MWIFI_ROUTING_TABLE_ADD:
    case MDF_EVENT_MWIFI_ROUTING_TABLE_REMOVE:
        MDF_LOGI("total_num: %d", esp_mesh_get_total_node_num());
        break;

    case MDF_EVENT_MWIFI_ROOT_GOT_IP:
    {
        MDF_LOGI("Root obtains the IP address. It is posted by LwIP stack automatically");
        xTaskCreate(tcp_client_write_task, "tcp_client_write_task", 4 * 1024,
                    NULL, CONFIG_MDF_TASK_DEFAULT_PRIOTY, NULL);
        xTaskCreate(tcp_client_read_task, "tcp_server_read", 4 * 1024,
                    NULL, CONFIG_MDF_TASK_DEFAULT_PRIOTY, NULL);
        break;
    }

    default:
        break;
    }

    return MDF_OK;
}

void setup_Weigand(__wIF *wif, uint8_t D0, uint8_t D1)
{

    //wif->gpio_evt_queue = xQueueCreate(40, sizeof(uint32_t));
    wif->_cardTempHigh = 0;
    wif->_cardTemp = 0;
    wif->_lastWiegand = 0;
    wif->_code = 0;
    wif->_bitCount = 0;
    wif->_wiegandType = 0;

    /*
	gpio_config_t io_conf;

    //interrupt of rising edge
    io_conf.intr_type = GPIO_PIN_INTR_NEGEDGE;
    //bit mask of the pins, use GPIO4/5 here
    io_conf.pin_bit_mask = (1ULL << D0) | (1ULL << D1);
    //set as input mode    
    io_conf.mode = GPIO_MODE_INPUT;
    //enable pull-up mode
    io_conf.pull_up_en = 1;
    gpio_config(&io_conf);
	*/
    gpio_config_t io_conf;

    //interrupt of rising edge
    io_conf.intr_type = GPIO_PIN_INTR_NEGEDGE;
    //bit mask of the pins, use GPIO4/5 here
    io_conf.pin_bit_mask = (/*GPIO_INPUT_PIN_SEL |*/ GPIO_INPUT_PIN_SEL2);
    //set as input mode
    io_conf.mode = GPIO_MODE_INPUT;
    //enable pull-up mode
    io_conf.pull_up_en = 1;
    gpio_config(&io_conf);

    //gpio_install_isr_service(ESP_INTR_FLAG_DEFAULT);

    //gpio_install_isr_service(ESP_INTR_FLAG_DEFAULT);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(GPIO_INPUT_IO_0, gpio_isr_handler, (void *)GPIO_INPUT_IO_0);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(GPIO_INPUT_IO_1, gpio_isr_handler, (void *)GPIO_INPUT_IO_1);
    /*
	gpio_isr_handler_add(GPIO_INPUT_IO_2, gpio_isr_handler, (void*) GPIO_INPUT_IO_2);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(GPIO_INPUT_IO_3, gpio_isr_handler, (void*) GPIO_INPUT_IO_3);
	 */
    /*
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(D0, gpio_isr_handler, (void*) D0);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(D1, gpio_isr_handler, (void*) D1);
	*/
    //xTaskCreate(gpio_task_example, "gpio_task_example", 2048, wif, 10, NULL);
    printf("Created %llu", wif->_cardTemp);
}

void app_main()
{
    mwifi_init_config_t cfg = MWIFI_INIT_CONFIG_DEFAULT();
    mwifi_config_t config = {
        .router_ssid = CONFIG_ROUTER_SSID,
        .router_password = CONFIG_ROUTER_PASSWORD,
        .mesh_id = CONFIG_MESH_ID,
        .mesh_password = CONFIG_MESH_PASSWORD,
    };

    /**
     * @brief Set the log level for serial port printing.
     */
    esp_log_level_set("*", ESP_LOG_INFO);
    esp_log_level_set(TAG, ESP_LOG_DEBUG);

    gpio_pad_select_gpio(CONFIG_LED_GPIO_NUM);
    gpio_set_direction(CONFIG_LED_GPIO_NUM, GPIO_MODE_INPUT_OUTPUT);

    /**
     * @brief Initialize wifi mesh.
     */
    MDF_ERROR_ASSERT(mdf_event_loop_init(event_loop_cb));
    MDF_ERROR_ASSERT(wifi_init());
    MDF_ERROR_ASSERT(mwifi_init(&cfg));
    MDF_ERROR_ASSERT(mwifi_set_config(&config));
    MDF_ERROR_ASSERT(mwifi_start());

    for (int i = 0; i < MAX_REPOST_DEPTH; i++)
    {
        cardMessageIds[i] = 0;
        cardRepostTime[i] = 0;
        cardIds[i] = 0;
    }

    gpio_config_t io_conf;

    //interrupt of rising edge
    io_conf.intr_type = GPIO_PIN_INTR_NEGEDGE;
    //bit mask of the pins, use GPIO4/5 here
    io_conf.pin_bit_mask = GPIO_INPUT_PIN_SEL | GPIO_INPUT_PIN_SEL2;
    //set as input mode
    io_conf.mode = GPIO_MODE_INPUT;
    //enable pull-up mode
    io_conf.pull_up_en = 1;
    gpio_config(&io_conf);

    //gpio_install_isr_service(ESP_INTR_FLAG_DEFAULT);

    __wIF *wif = &(wiegandTerm[0]);

    wif->_cardTempHigh = 0;
    wif->_cardTemp = 0;
    wif->_lastWiegand = 0;
    wif->_code = 0;
    wif->_cardRead = 0;
    wif->_bitCount = 0;
    wif->_wiegandType = 0;

    xTaskCreate(gpio_task_example, "gpio_task_example", 2048, wif, 10, NULL);

    wif = &(wiegandTerm[1]);

    wif->_cardTempHigh = 0;
    wif->_cardTemp = 0;
    wif->_lastWiegand = 0;
    wif->_code = 0;
    wif->_cardRead = 0;
    wif->_bitCount = 0;
    wif->_wiegandType = 0;

    xTaskCreate(gpio_task_example, "gpio_task_example2", 2048, wif, 10, NULL);

    messageID = 1;
    //install gpio isr service
    gpio_install_isr_service(ESP_INTR_FLAG_DEFAULT);
    //hook isr handler for specific gpio pin

    gpio_isr_handler_add(GPIO_INPUT_IO_0, gpio_isr_handler, (void*) GPIO_INPUT_IO_0);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(GPIO_INPUT_IO_1, gpio_isr_handler, (void*) GPIO_INPUT_IO_1);
 
    gpio_isr_handler_add(GPIO_INPUT_IO_2, gpio_isr_handler, (void *)GPIO_INPUT_IO_2);
    //hook isr handler for specific gpio pin
    gpio_isr_handler_add(GPIO_INPUT_IO_3, gpio_isr_handler, (void *)GPIO_INPUT_IO_3);

    // CREATE RFID ACCESS QUEUE

    RFID_ACCESS_Q = xQueueCreate(10, sizeof(unsigned long long));
    TAG_ACCESS_Q = xQueueCreate(10, sizeof(unsigned long long));

    xTaskCreate(gateway_heartbeat, "gateway_heartbeat", 2048, NULL, 10, NULL);


    xTaskCreate(checkPendingID, "checkPendingID", 2048, NULL, 10, NULL);


    /**
     * @breif Create handler
     */
    xTaskCreate(node_write_task, "node_write_task", 4 * 1024,
                NULL, CONFIG_MDF_TASK_DEFAULT_PRIOTY, NULL);
    xTaskCreate(node_read_task, "node_read_task", 4 * 1024,
                NULL, CONFIG_MDF_TASK_DEFAULT_PRIOTY, NULL);

    TimerHandle_t timer = xTimerCreate("print_system_info", 10000 / portTICK_RATE_MS,
                                       true, NULL, print_system_info_timercb);
    xTimerStart(timer, 0);

    uart_init();
    xTaskCreate(uart_rx_task, "uart_rx_task", 1024 * 2, NULL, 10, NULL);



}
