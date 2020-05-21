package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const plugins = `
CREATE EXTENSION IF NOT EXISTS unaccent;
`

const importableTypes = `
CREATE TYPE importable AS ENUM ('Properties');
`

const cardLinkTypes = `
CREATE TYPE cardlink AS ENUM ('Account', 'Property');
`

const cardTypes = `
CREATE TYPE cardtype AS ENUM ('EM4100', 'MIFARE 1K', 'MIFARE 4K', 'HID Card', 'T5567', '2nd Card', 'ISO14443B', 'FELICA', '15693 Label', 'CPU Card');
`

const transitType = `
CREATE TYPE transittype AS ENUM ('ENTER', 'EXIT', 'POSSES');
`

const accountTable = `
CREATE TABLE accounts (
id serial NOT NULL,
uuid uuid NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone DEFAULT current_timestamp,
last_login timestamp with time zone NOT NULL DEFAULT current_timestamp,
email text NOT NULL UNIQUE,
pass text NOT NULL,
name text NOT NULL,
active boolean NOT NULL DEFAULT TRUE,
roles text[] NOT NULL DEFAULT '{"user"}',
card_id bigint,
PRIMARY KEY (id),
UNIQUE(uuid),
    CONSTRAINT card FOREIGN KEY (card_id)
        REFERENCES public.cards (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)`

const tokenTable = `
CREATE TABLE tokens (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
account_id uuid NOT NULL,
token text NOT NULL UNIQUE,
expiry timestamp with time zone NOT NULL,
mobile boolean NOT NULL DEFAULT FALSE,
identifier text,
PRIMARY KEY (id),
CONSTRAINT accounts FOREIGN KEY (account_id)
REFERENCES public.accounts (uuid) MATCH SIMPLE
ON UPDATE NO ACTION
ON DELETE NO ACTION
)`

const doorAccessTable = `
CREATE TABLE door_accesses
(
    mac macaddr,
	id serial,
	enter_at timestamp without time zone,
	exit_at timestamp without time zone,
	exit_type smallint,
	card_code text,
    card_id bigint,
    gateway_id bigint,
    PRIMARY KEY (id),
	CONSTRAINT card FOREIGN KEY (card_id)
        REFERENCES public.cards (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
	CONSTRAINT gateway FOREIGN KEY (gateway_id)
			REFERENCES public.gateways (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
)
    
`

const cardTable = `
CREATE TABLE cards
(
	code text NOT NULL,
	uuid uuid NOT NULL,
    id serial NOT NULL,
	type cardtype,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	link cardlink,
	link_id integer,
    CONSTRAINT card_pkey PRIMARY KEY (id)
)
`

const cardsLogTable = `
CREATE TABLE card_logs
(
	code text NOT NULL,
    id serial NOT NULL,
	type cardtype,
    gateway_id bigint,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	CONSTRAINT card_log_pkey PRIMARY KEY (id),
	CONSTRAINT gateway FOREIGN KEY (gateway_id)
			REFERENCES public.gateways (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
)
`

const roomTable = `
CREATE TABLE rooms
(
    id serial NOT NULL,
	uuid uuid NOT NULL,
	name varchar(60) NOT NULL,
	code varchar(60),
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    active_user_id bigint,
	CONSTRAINT room_pkey PRIMARY KEY (id),
	CONSTRAINT active_user FOREIGN KEY (active_user_id)
		REFERENCES public.accounts (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION
)
`

const gatewayTable = `
CREATE TABLE gateways
(
    id serial NOT NULL,
	uuid uuid NOT NULL,
	name varchar(60),
	addr macaddr UNIQUE NOT NULL,
	room_id bigint references rooms(id),
	active boolean DEFAULT 'f',
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	last_active timestamp with time zone NOT NULL DEFAULT current_timestamp,
	CONSTRAINT gateway_pkey PRIMARY KEY (id),
	CONSTRAINT room FOREIGN KEY (room_id)
		REFERENCES public.rooms (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION
)
`

const filesTable = `
CREATE TABLE public.files
(
    id serial NOT NULL,
    uuid uuid,
    mime varchar(127),
    name varchar(256),
    size bigint,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	account_id bigint,
    CONSTRAINT file_pkey PRIMARY KEY (id),
	CONSTRAINT account FOREIGN KEY (account_id)
		REFERENCES public.accounts (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION
)
`

const propertyTable = `
CREATE TABLE properties
(
    id serial NOT NULL,
	uuid uuid NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	name varchar(256) NOT NULL,
	code varchar(128),
	price real,
	acquired_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    room_id bigint,
	CONSTRAINT property_pkey PRIMARY KEY (id),
	CONSTRAINT room FOREIGN KEY (room_id)
		REFERENCES public.rooms (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION
)
`

const propertyTransitTable = `
CREATE TABLE property_transits
(
    id serial NOT NULL,
    card_id bigint,
    property_id bigint,
    room_id bigint,
    user_id bigint,
	created_at timestamp without time zone,
	transition transittype,
	comment text,
	post_auth boolean,
	PRIMARY KEY (id),
	CONSTRAINT properties FOREIGN KEY (property_id)
		REFERENCES public.properties (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION,
	CONSTRAINT cards FOREIGN KEY (card_id)
		REFERENCES public.cards (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION,
	CONSTRAINT room FOREIGN KEY (room_id)
		REFERENCES public.rooms (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION,
	CONSTRAINT active_user FOREIGN KEY (user_id)
		REFERENCES public.accounts (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION
)
`

const importTable = `
CREATE TABLE imports
(
    id serial NOT NULL,
	uuid uuid NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	created_by_id integer NOT NULL,
	resource importable NOT NULL,
	elements integer[],
	
	CONSTRAINT import_pkey PRIMARY KEY (id)
)
`

func init() {
	up := []string{
		plugins,
		importableTypes,
		cardLinkTypes,
		cardTypes,
		cardTable,
		transitType,
		accountTable,
		tokenTable,
		roomTable,
		gatewayTable,
		doorAccessTable,
		filesTable,
		propertyTable,
		propertyTransitTable,
		importTable,
		cardsLogTable,
	}

	down := []string{
		`DROP TABLE tokens`,
		`DROP TABLE accounts`,
		`DROP TABLE door_access`,
		`DROP TABLE cards`,
		`DROP TABLE gateways`,
		`DROP TABLE rooms`,
		`DROP TABLE files`,
		`DROP TABLE properties`,
		`DROP TABLE imports`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
