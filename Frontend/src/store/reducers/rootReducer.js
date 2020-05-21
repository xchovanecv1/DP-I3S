import { combineReducers } from 'redux';
import authReducer from './authReducer';
import accountReducer from './accountReducer';
import roomReducer from './roomReducer';

export default combineReducers({
    authReducer,
    accountReducer,
    roomReducer
});