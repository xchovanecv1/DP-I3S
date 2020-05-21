import axios from "../../../common/axios"
import { thunk, action } from 'easy-peasy';

import schema from "./schema";
import merge from "deepmerge";

function formatError(data) {
	if (!!!data) return null;
	return {
		code: data.status,
		codeError: data.statusText,
		text: data.data && data.data.status ? data.data.status : null,
	}
}

export const model = {
	props: {
		properties: [],
		count: 0,
	},
	isLoading: false,
	error: null,

	prop: schema,
	propLoading: false,
	propError: null,

	initalLoaded: false,
	setInitialLoaded: action((state) => {
		if (!state.initalLoaded) {
			state.initalLoaded = true;
		}
	}),

	setProp: action((state, payload) => {
		state.prop = payload;
	}),
	setPropLoading: action((state, payload) => {
		state.propLoading = payload;
	}),
	setPropError: action((state, payload) => {
		state.propError = formatError(payload);
	}),
	loadProp: thunk(async (actions, payload) => {
		actions.setPropLoading(true);
		actions.setProp(actions.getSchema());
		axios.get(`/api/admin/props/${payload}`)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				const data = merge(schema, res.data);
				actions.setProp(data);
			}
		})
		.catch(err => {
			actions.setPropError(err.response);
		})
		.finally(() => {
			actions.setPropLoading(false);
		});
	}),
	saveProp: thunk(async (actions, payload) => {
		actions.setPropLoading(true);
		//actions.setProp(actions.getSchema());
		axios.put(`/api/admin/props/${payload.id}`, payload)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				const data = merge(schema, res.data);
				actions.setProp(data);
			}
		})
		.catch(err => {
			actions.setPropError(err.response);
		})
		.finally(() => {
			actions.setPropLoading(false);
		});
	}),
	addProp: thunk(async (actions, payload) => {
		actions.setPropLoading(true);
		//actions.setProp(actions.getSchema());
		return axios.post(`/api/admin/props/`, payload)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				const data = merge(schema, res.data);
				actions.setProp(data);
				return data;
			}
		})
		.catch(err => {
			actions.setPropError(err.response);
		})
		.finally(() => {
			actions.setPropLoading(false);
		});
	}),
	getSchema: thunk(async () => {
		return merge({}, schema);
	}),

	clearProp: thunk(async (actions) => {
		actions.setProp(merge({}, schema));
	}),

	// IMPORT 
	importFile: null,

	headerLoading: false,
	headerError: null,
	importHeader: null,

	fileLoading: false,
	fileError: null,
	fileHeader: null,

	importLoading: false,
	importError: null,
	importReturn: null,

	commitPending: false,
	commitError: null,

	setImportFile: action((state, payload) => {
		state.importFile = payload;
	}),

	setHeaderLoading: action((state, payload) => {
		state.headerLoading = payload;
	}),
	setHeaderError: action((state, payload) => {
		state.headerError = payload;
	}),
	setImportHeader: action((state, payload) => {
		state.importHeader = payload;
	}),

	setFileLoading: action((state, payload) => {
		state.fileLoading = payload;
	}),
	setFileError: action((state, payload) => {
		state.fileError = payload;
	}),
	setFileHeader: action((state, payload) => {
		state.fileHeader = payload;
	}),

	setImportLoading: action((state, payload) => {
		state.importLoading = payload;
	}),
	setImportError: action((state, payload) => {
		state.importError = payload;
	}),
	setImportReturn: action((state, payload) => {
		state.importReturn = payload;
	}),


	setCommitError: action((state, payload) => {
		state.commitError = payload;
	}),
	setCommitPending: action((state, payload) => {
		state.commitPending = payload;
	}),

	clearImport: thunk(async (actions, payload) => {
		actions.setImportFile(null);
		actions.setHeaderLoading(false);
		actions.setHeaderError(null);
		actions.setImportHeader(null);
		actions.setFileLoading(false);
		actions.setFileError(null);
		actions.setFileHeader(null);
		actions.setImportLoading(false);
		actions.setImportError(null);
		actions.setImportReturn(null);
	}),

	loadImportHeader: thunk(async (actions, payload) => {
		actions.setHeaderLoading(true);
		const ret = axios.get("/api/admin/props/import/headers")
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setImportHeader(res.data);
				actions.setHeaderLoading(false);
			}
		})
		.catch(err => {
			actions.setHeaderLoading(false);
			actions.setHeaderError(err.response);
		});
	}),
	loadFileHeader: thunk(async (actions, payload) => {
		actions.setFileLoading(true);
		const ret = axios.get(`/api/admin/files/${payload}/head`)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setFileHeader(res.data);
				actions.setFileLoading(false);
			}
		})
		.catch(err => {
			actions.setFileLoading(false);
			actions.setFileError(err.response);
		});
	}),

	importFile: thunk(async (actions, payload) => {
		actions.setImportLoading(true);
		let url = `/api/admin/props/import/${payload.id}`;
		if (payload.head) {
			url += "?header=t";
		}
		const ret = axios.post(url, payload.map)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setImportReturn(res.data);
				actions.setImportLoading(false);
				return res.data;
			}
		})
		.catch(err => {
			actions.setImportLoading(false);
			actions.setImportError(err.response);
		});
		actions.setImportReturn(ret);
	}),

	commitImport: thunk(async (actions, payload) => {
		actions.setCommitPending(true);
		let url = `/api/admin/props/commit`;
		const ret = axios.post(url, payload)
		.then(() => {
				actions.setCommitPending(false);
				actions.setCommitError(null);
		})
		.catch(err => {
			actions.setCommitPending(false);
			actions.setCommitError(err.response);
		})
		.finally(() => {
			actions.clearImport();
		});
		return ret;
	}),
	//
	setLoading: action((state, payload) => {
		state.isLoading = payload;
	}),
	setError: action((state, payload) => {
		state.error = payload;
	}),
	setProps: action((state, payload) => {
		state.props = payload;
		state.isLoading = false;
		state.error = null;
	}),
	loadAll: thunk(async (actions, payload) => {
		actions.setLoading(true);
		let url = `/api/admin/props`
		if (payload) {
			url = `${url}?q=${encodeURIComponent(payload)}`;
		}
		axios.get(url)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setProps(res.data);
				actions.setInitialLoaded();
			}
		})
		.catch(err => {
			actions.setError(err.response);
		})
		.finally(() => {
			actions.setLoading(false);
		});
	}),


	///
	location: null,
	locationLoading: false,
	locationError: null,

	setLocation: action((state, payload) => {
		state.location = payload;
	}),
	setLocationLoading: action((state, payload) => {
		state.locationLoading = payload;
	}),
	setLocationError: action((state, payload) => {
		state.locationError = formatError(payload);
	}),
	loadLocation: thunk(async (actions, payload) => {
		actions.setLocationLoading(true);
		axios.get(`/api/admin/props/${payload}/location`)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setLocation(res.data);
			}
		})
		.catch(err => {
			actions.setLocationError(err.response);
		})
		.finally(() => {
			actions.setLocationLoading(false);
		});
	}),


	///
	possesPending: [],

	setPossessPending: action((state, payload) => {
		state.possesPending = (payload);
	}),
	loadPendinProps: thunk(async (actions, payload) => {
		return axios.get(`/api/admin/props/pending`)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setPossessPending(res.data);
			}
		})
	}),

	savePendinProps: thunk(async (actions, payload) => {
		return axios.post(`/api/admin/props/pending`, payload)
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setPossessPending([]);
			}
		})
	}),

	///

	//
	cardLinkDeleting: false,
	cardLinkDeleteError: null,

	setCardLinkDeleting: action((state, payload) => {
		state.cardLinkDeleting = payload;
	}),
	setCardLinkDeleteError: action((state, payload) => {
		state.cardLinkDeleteError = formatError(payload);
	}),
	deleteCardLink: thunk(async (actions, payload) => {
		actions.setCardLinkDeleting(true);
		return axios.delete(`/api/admin/props/${payload.id}/cards`, {
			data: payload.data
		})
			.then(() => {
				return true;
			})
			.catch(err => {
				actions.setCardLinkDeleteError(err.response);
				return false;
			})
			.finally(() => {
				actions.setCardLinkDeleting(false);
			});
	}),

	//
	cardLinkAdding: false,
	cardLinkAddingError: null,

	setCardLinkAdding: action((state, payload) => {
		state.cardLinkAdding = payload;
	}),
	setCardLinkAddingError: action((state, payload) => {
		state.cardLinkAddingError = formatError(payload);
	}),
	addCardLinkFromLog: thunk(async (actions, payload) => {
		actions.setCardLinkAdding(true);
		return axios.post(`/api/admin/props/${payload.id}/cards/log`, payload.data)
			.then(() => {
				return true;
			})
			.catch(err => {
				actions.setCardLinkAddingError(err.response);
				return false;
			})
			.finally(() => {
				actions.setCardLinkAdding(false);
			});
	}),
}