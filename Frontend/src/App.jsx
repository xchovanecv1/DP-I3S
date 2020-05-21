import React, { Component, useEffect, useState } from 'react';
import { Router, Route, Redirect } from 'react-router-dom';
import { createBrowserHistory } from 'history';

import { connect } from 'react-redux';

import { useStoreState, useStore, useStoreActions } from 'easy-peasy'

// Externals
import { Chart } from 'react-chartjs-2';

// Material helpers
import { ThemeProvider } from '@material-ui/styles';

// ChartJS helpers
import { chartjs } from './helpers';

// Theme
import theme from './theme';

// Styles
import 'react-perfect-scrollbar/dist/css/styles.css';
import './assets/scss/index.scss';

import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import MomentUtils from '@date-io/date-fns';

// Routes
import Routes from './Routes';


// Browser history
const browserHistory = createBrowserHistory();

// Configure ChartJS
Chart.helpers.extend(Chart.elements.Rectangle.prototype, {
  draw: chartjs.draw
});

browserHistory.listen((location, action) => {
  console.log("[ROUTER]",location, action);
});

window.browserHistory = browserHistory;


const App = (props) => {

  const isLogged = useStoreState(state => state.auth.authenticated);
  const initTokens = useStoreActions(actions => actions.auth.initTokens);
  const [ inited, setInited ] = useState(false);
  useEffect(() => {
      initTokens();
      setInited(true);
  }, []);

  return (
      <ThemeProvider theme={theme}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          {inited && 
            <Router history={browserHistory}>
              <Routes isLogged={isLogged} />
            </Router>
          }
        </MuiPickersUtilsProvider>
      </ThemeProvider>
  );
}

export default (App);