import React from 'react';
import ReactDOM from 'react-dom';


import { StoreProvider } from 'easy-peasy';
//import { Provider } from 'react-redux'
import stateStore from './store';

// Service worker
import * as serviceWorker from './common/serviceWorker';

// App
import App from './App';

ReactDOM.render(
    <StoreProvider store={stateStore()}>
     <App />
    </StoreProvider>,
    document.getElementById('root')
);

serviceWorker.unregister();
