import axios from "common/axios"
import merge from "deepmerge"
import schema from "./schema"
import { thunk, action } from 'easy-peasy';

function formatError(data) {
	if (!!!data) return null;
	return {
		code: data.status,
		codeError: data.statusText,
		text: data.data && data.data.status ? data.data.status : null,
	}
}

export const model = {
	cards: {
		cards: [],
	},
	isLoading: false,
	error: null,

	card: schema,
	cardLoading: false,
	cardError: null,

	setCard: action((state, payload) => {
		state.card = payload;
	}),
	setCardLoading: action((state, payload) => {
		state.cardLoading = payload;
	}),
	setCardError: action((state, payload) => {
		state.cardError = formatError(payload);
	}),
	loadCard: thunk(async (actions, payload) => {
		actions.setCardLoading(true);
		actions.setCard(actions.getSchema());
		axios.get(`/api/admin/${payload.ctx}/${payload.id}/cards`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					const data = merge(schema, res.data);
					actions.setCard(data);
				}
			})
			.catch(err => {
				actions.setCardError(err.response);
			})
			.finally(() => {
				actions.setCardLoading(false);
			});
	}),
	getSchema: thunk(async () => {
		return merge({}, schema);
	}),

	logSearch: {
		cards: [],
	},
	logSearchLoading: false,
	logSearchError: null,

	setLogSearch: action((state, payload) => {
		state.logSearch = payload;
	}),
	setLogSearchLoading: action((state, payload) => {
		state.logSearchLoading = payload;
	}),
	setLogSearchError: action((state, payload) => {
		state.logSearchError = formatError(payload);
	}),
	searchCardLogs: thunk(async (actions, payload) => {
		actions.setLogSearchLoading(true);
		axios.get(`/api/admin/cards/logs?q=${encodeURIComponent(payload)}`)
			.then(res => {
				if (
					typeof res.data !== "undefined"
				) {
					actions.setLogSearch(res.data);
				}
			})
			.catch(err => {
				actions.setLogSearchError(err.response);
			})
			.finally(() => {
				actions.setLogSearchLoading(false);
			});
	}),

}