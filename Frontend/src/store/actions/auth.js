import axios from "../../common/axios"
import { thunk, action, computed } from 'easy-peasy';

let refresh_timer = null;

//https://github.com/marosivanco/keyclops/
function decodePayload(token) {
	function toUtf8Char(str) {
		try {
			return decodeURIComponent(str);
		} catch (err) {
			return String.fromCharCode(0xfffd); // UTF 8 invalid char
		}
	}
	function toUtf8String(buf) {
		let res = "";
		let tmp = "";
		for (var i = 0; i < buf.length; i++) {
			if (buf[i] <= 0x7f) {
				res = `${res}${toUtf8Char(tmp)}${String.fromCharCode(buf[i])}`;
				tmp = "";
			} else {
				tmp = `${tmp}%${buf[i].toString(16)}`;
			}
		}
		return res + toUtf8Char(tmp);
	}
	const base = token
		.split(".")[1]
		.replace(/-/g, "+")
		.replace(/_/g, "/");
	const data = atob(base);
	const bytes = new Uint8Array(new ArrayBuffer(data.length));
	for (var i = 0; i < data.length; i++) {
		bytes[i] = data.charCodeAt(i);
	}
	return JSON.parse(toUtf8String(bytes));
}

export const logout = () => {
	delete axios.defaults.headers.common["Authorization"];
	localStorage.removeItem("i3s-at");
	localStorage.removeItem("i3s-rf");
}

export const initTokens = (setTokens = () =>{}, setToken = () =>{}) => {
	const at = localStorage.getItem("i3s-at");
	const rt = localStorage.getItem("i3s-rf");
	console.log("[AUTH]", "Init");
	if(at && rt) {
		clearTimeout(refresh_timer);
		const token = decodePayload(at);

		const now = Date.now()/1000;
    	const ref_p = (token.exp- now - 60)*100;
		if(now > token.exp) {
			localStorage.removeItem("i3s-at");
			localStorage.removeItem("i3s-rf");
			console.log("[AUTH]", "Invalid token, logout");
			return false;
		}

		axios.defaults.headers.common["Authorization"] = `Bearer ${at}`;

		setTokens({
			access_token: at,
			refresh_token: rt,
		});
		setToken(token);

		refresh_timer = setTimeout(() => refreshToken(rt), ref_p);
		console.log("[AUTH]", "Next refesh in ", ref_p / 60000);

		return true;
	}
}

const updateTokens = (acc_t, ref_t, success = (e) => e) => {

		clearTimeout(refresh_timer);
    const token = decodePayload(acc_t);

    const now = Date.now()/1000;
    const ref_p = (token.exp- now - 60)*100;

		axios.defaults.headers.common["Authorization"] = `Bearer ${acc_t}`;

    localStorage.setItem('i3s-at', acc_t);
    localStorage.setItem('i3s-rf', ref_t);

		refresh_timer = setTimeout(() => refreshToken(ref_t), ref_p);
		console.log("[AUTH]", "Next refesh in ", ref_p / 60000);
		success(token);
	return token;

}

const refreshToken = (ref_token, success = (e) => e) => {

    delete axios.defaults.headers.common["Authorization"];
    const ret = axios.post("/api/auth/refresh", null, {
        headers: {
            Authorization: `Bearer ${ref_token}`,
        }
    })
    .then(res => { 
			if (
					typeof res.data !== "undefined" &&
					typeof res.data.access_token !== "undefined" &&
					typeof res.data.refresh_token !== "undefined"
			) {
					const acc_t = res.data.access_token;
					const ref_t = res.data.refresh_token;
					
					return updateTokens(acc_t, ref_t, success);
			}
    });

    return ret;
}

export const testApi = () => {
	return axios.get("/api/admin");
}

export const loginUser = (login, password, success = () => {}) => {
		const ret = axios.post("/api/auth/login", {
			Email: login,
			Pass: password,
		})
		.then(res => {
			if (
				typeof res.data !== "undefined" &&
				typeof res.data.access_token !== "undefined" &&
				typeof res.data.refresh_token !== "undefined"
			) {
				const acc_t = res.data.access_token;
				const ref_t = res.data.refresh_token;
				
				return updateTokens(acc_t, ref_t, success);
			}
		});
		return ret;
		//return axios.post("/api/auth/login", { email: login, pass: password });
}

export const REQUEST_LOGIN = 'REQUEST_LOGIN'
function requestLogin(name) {
  return {
    type: REQUEST_LOGIN,
    name
  }
}

export const LOGIN_SUCCESSFUL = 'LOGIN_SUCCESSFUL'
export function loginSuccessful(token) {
  return {
    type: LOGIN_SUCCESSFUL,
    token,
    loggedAt: Date.now()
  }
}

export const model = {
	access_token: null,
	refresh_token: null,
	token: {},
	isLoading: false,
	error: null,
	authenticated: computed(state => state.access_token !== null),
	setLoading: action((state, payload) => {
		state.isLoading = payload;
	}),
	setError: action((state, payload) => {
		state.error = payload;
	}),
	setToken: action((state, payload) => {
		state.token = {
			...payload
		};
	}),
	setLoginResponse: action((state, payload) => {
		state.access_token = payload.access_token;
		state.refresh_token = payload.refresh_token;
	}),
	initTokens: thunk(async (actions, payload) => {
		return initTokens(actions.setLoginResponse, actions.setToken);
	}),
	logout: action((state, payload) => {
		state.access_token = null;
		state.refresh_token = null;
		state.token = {};
		logout();
	}),
	login: thunk(async (actions, payload) => {
		actions.setLoading(true);
		const ret = axios.post("/api/auth/login", {
			Email: payload.email,
			Pass: payload.password,
		})
		.then(res => {
			if (
				typeof res.data !== "undefined" &&
				typeof res.data.access_token !== "undefined" &&
				typeof res.data.refresh_token !== "undefined"
			) {
				const acc_t = res.data.access_token;
				const ref_t = res.data.refresh_token;
				actions.setLoginResponse({
					...res.data
				});
				const token = updateTokens(acc_t, ref_t);
				actions.setToken(token);
				actions.setLoading(false);
			}
		})
		.catch(err => {
			actions.setLoading(false);
			actions.setError(err.response);
		})
		;
		/*const updated = await loginUser.update(payload.id, payload);
		actions.setProduct(updated); // 👈 dispatch local actions to update state*/
	}),
}