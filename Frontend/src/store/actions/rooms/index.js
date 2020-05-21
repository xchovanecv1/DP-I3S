import axios from "common/axios"
import merge from "deepmerge"
import schema from "./schema"
import { thunk, action } from 'easy-peasy';

function formatError(data) {
	if (!!!data) return null;
	console.log(data);
	return {
		code: data.status,
		codeError: data.statusText,
		text: data.data && data.data.status ? data.data.status : null,
	}
}

export const model = {
	rooms: {
		properties: [],
		count: 0,
	},
	isLoading: false,
	error: null,
	firstLoadFinished: false,

	room: schema,
	roomLoading: false,
	roomError: null,

	setRoom: action((state, payload) => {
		state.room = payload;
	}),
	setRoomLoading: action((state, payload) => {
		state.roomLoading = payload;
	}),
	setRoomError: action((state, payload) => {
		state.roomError = formatError(payload);
	}),
	loadRoom: thunk(async (actions, payload) => {
		actions.setRoomLoading(true);
		actions.setRoom(actions.getSchema());
		axios.get(`/api/admin/rooms/${payload}`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					const data = merge(schema, res.data);
					actions.setRoom(data);
				}
			})
			.catch(err => {
				actions.setRoomError(err.response);
			})
			.finally(() => {
				actions.setRoomLoading(false);
			});
	}),
	getSchema: thunk(async () => {
		return merge({}, schema);
	}),

	search: {
		properties: [],
		count: 0,
	},
	searchLoading: false,
	searchError: null,

	setSearch: action((state, payload) => {
		state.search = payload;
	}),
	setSearchLoading: action((state, payload) => {
		state.searchLoading = payload;
	}),
	setSearchError: action((state, payload) => {
		state.searchError = formatError(payload);
	}),
	searchRooms: thunk(async (actions, payload) => {
		actions.setSearchLoading(true);
		axios.get(`/api/admin/rooms?q=${encodeURIComponent(payload)}`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					actions.setSearch(res.data);
				}
			})
			.catch(err => {
				actions.setSearchError(err.response);
			})
			.finally(() => {
				actions.setSearchLoading(false);
			});
	}),

	setLoading: action((state, payload) => {
		state.isLoading = payload;
	}),
	setError: action((state, payload) => {
		state.error = payload;
	}),
	setRooms: action((state, payload) => {
		state.rooms = payload;
		state.isLoading = false;
		state.error = null;
	}),
	setFirstLoad: action((state, payload) => {
		if (!state.firstLoadFinished) {
			state.firstLoadFinished = true;
		}
	}),
	loadAll: thunk(async (actions, payload) => {
		let url = `/api/admin/rooms`
		if (payload) {
			actions.setSearchLoading(true);
			url = `${url}?q=${encodeURIComponent(payload)}`;
		} else {
			actions.setLoading(true);
		}
		axios.get(url)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setRooms(res.data);
			}
		})
		.catch(err => {
			actions.setError(err.response);
		})
		.finally(() => {
			actions.setLoading(false);
			actions.setSearchLoading(false);
			actions.setFirstLoad();
		});
	}),

	saving: false,
	saveError: null,

	setSaving: action((state, payload) => {
		state.saving = payload;
	}),
	setSaveError: action((state, payload) => {
		state.saveError = payload;
	}),

	saveRoom: thunk(async (actions, payload) => {
		actions.setLoading(true);
		//actions.setProp(actions.getSchema());
		axios.put(`/api/admin/rooms/${payload.id}`, payload)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				const data = merge(schema, res.data);
				actions.setRoom(data);
			}
		})
		.catch(err => {
			actions.setRoomError(err.response);
		})
		.finally(() => {
			actions.setRoomLoading(false);
		});
	}),

	addRoom: thunk(async (actions, payload) => {
		actions.setLoading(true);
		//actions.setProp(actions.getSchema());
		return axios.post(`/api/admin/rooms`, payload)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				const data = merge(schema, res.data);
				actions.setRoom(data);
				return data;
			}
		})
		.catch(err => {
			actions.setRoomError(err.response);
		})
		.finally(() => {
			actions.setRoomLoading(false);
		});
	}),

	clearRoom: thunk(async (actions) => {
		actions.setRoom(merge({}, schema));
	}),

	gateways: [],
	gatewaysLoading: false,
	gatewaysError: null,

	setGateways: action((state, payload) => {
		state.gateways = payload;
	}),
	setGatewaysLoading: action((state, payload) => {
		state.gatewaysLoading = payload;
	}),
	setGatewaysError: action((state, payload) => {
		state.gatewaysError = formatError(payload);
	}),
	loadGateways: thunk(async (actions, payload) => {
		actions.setGatewaysLoading(true);
		axios.get(`/api/admin/rooms/${encodeURIComponent(payload)}/gateway`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					actions.setGateways(res.data);
				}
			})
			.catch(err => {
				actions.setGatewaysError(err.response);
			})
			.finally(() => {
				actions.setGatewaysLoading(false);
			});
	}),


	gwAddPending: false,
	gwAddError: null,

	setGwAddPending: action((state, payload) => {
		state.gwAddPending = payload;
	}),
	setGwAddError: action((state, payload) => {
		state.gwAddError = payload;
	}),
	addRoomGw: thunk(async (actions, payload) => {
		actions.setGwAddPending(true);
		//actions.setProp(actions.getSchema());
		axios.post(`/api/admin/rooms/${payload.id}/gateway`, payload.gateway)
		.catch(err => {
			actions.setGwAddError(err.response);
		})
		.finally(() => {
			actions.setGwAddPending(false);
		});
	}),


	gws: {
		gateways: [],
	},
	gwSearchLoading: false,
	gwSearchError: null,

	setGwSearch: action((state, payload) => {
		state.gws = payload;
	}),
	setGwSearchLoading: action((state, payload) => {
		state.gwSearchLoading = payload;
	}),
	setGwSearchError: action((state, payload) => {
		state.logSearchError = formatError(payload);
	}),
	searchGateways: thunk(async (actions, payload) => {
		actions.setGwSearchLoading(true);
		axios.get(`/api/admin/gateways/unused?q=${encodeURIComponent(payload)}`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					actions.setGwSearch(res.data);
				}
			})
			.catch(err => {
				actions.setGwSearchError(err.response);
			})
			.finally(() => {
				actions.setGwSearchLoading(false);
			});
	}),
}