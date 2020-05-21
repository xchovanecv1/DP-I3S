import axios from "../../common/axios"
import { thunk, action, computed } from 'easy-peasy';


function formatError(data) {
	if (!!!data) return null;
	return {
		code: data.status,
		codeError: data.statusText,
		text: data.data && data.data.status ? data.data.status : null,
	}
}


export const ALL_ACCOUNTS_LOADED = 'ALL_ACCOUNTS_LOADED'
export function allAccountsLoaded(accounts) {
  return {
    type: ALL_ACCOUNTS_LOADED,
    accountsPending: false,
    accounts,
  }
}

export const ALL_ACCOUNTS_PENDING = 'ALL_ACCOUNTS_PENDING'
export function allAccountsPending() {
  return {
    type: ALL_ACCOUNTS_PENDING,
    accountsPending: true,
  }
}

export const ALL_ACCOUNTS_FAILED = 'ALL_ACCOUNTS_FAILED'
export function allAccountsFailed(error) {
  return {
    type: ALL_ACCOUNTS_FAILED,
    accountsPending: false,
    accountsError: error,
  }
}

export function loadAllAccounts(page, perPage = 10) {
    return dispatch => {
        dispatch(allAccountsPending());
        const url = `/api/admin/accounts?limit=${perPage ? perPage : 1}&page=${page}`;
        axios.get(url)
        .then(res => {
            dispatch(allAccountsLoaded(res.data));
        })
        .catch(err => {
            dispatch(allAccountsFailed(err.response));
        });
    }
}


export const model = {
	accounts: [],
	isLoading: false,
	error: null,
	setLoading: action((state, payload) => {
		state.isLoading = payload;
	}),
	setError: action((state, payload) => {
		state.error = payload;
  }),
  setAccounts: action((state, payload) => {
    state.accounts = payload
  }),
	loadAll: thunk(async (actions, payload) => {
    const { page, perPage } = payload;
		actions.setLoading(true);
		const url = `/api/admin/accounts?limit=${perPage ? perPage : 1}&page=${page}`;
    axios.get(url)
    .then(res => {
      actions.setAccounts(res.data);
    })
    .catch(err => {
      actions.setError(err.response);
    });
  }),
  
  //

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
		return axios.delete(`/api/admin/accounts/${payload.id}/cards`, {
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
		return axios.post(`/api/admin/accounts/${payload.id}/cards/log`, payload.data)
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