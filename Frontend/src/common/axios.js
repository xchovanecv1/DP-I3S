import axios from "axios";

const Axios = axios.create({
	timeout: 15000,
});

export default Axios;