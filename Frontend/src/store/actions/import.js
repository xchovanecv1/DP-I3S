import axios from "../../common/axios"
import { thunk, action, computed } from 'easy-peasy';

export const model = {

    actionList: [],
    actionError: null,
    actionLoading: false,

    setActionList: action((state, payload) => {
		state.actionList = payload;
	}),
	setActionError: action((state, payload) => {
		state.actionError = payload;
	}),
	setActionLodaing: action((state, payload) => {
		state.actionLoading = payload;
    }),
    
    loadActionList: thunk(async (actions, payload) => {
		actions.setActionLodaing(true);
		const ret = axios.get("/api/admin/import/actions")
		.then(res => {
			if (
				typeof res.data !== "undefined"
			) {
				actions.setActionList(res.data);
				actions.setActionLodaing(false);
			}
		})
		.catch(err => {
			actions.setActionLodaing(false);
			actions.setActionError(err.response);
		});
	}),
    
	upload: thunk(async (actions, payload) => {
        
        const data = new FormData()
        data.append('file', payload)

        const ret = axios.post("/api/admin/files/upload", data, { 
            'Content-Type': 'application/x-www-form-urlencoded'
        }).then(r => {
            if(r.data && r.data.id) {
                return r.data;
            } else {
                return new Promise((res, rej) => {
                    rej(false);
                });
            }
        })

        return ret;

	}),
}