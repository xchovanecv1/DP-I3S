import React, { Component, useState } from 'react';
import { Link, withRouter } from 'react-router-dom';

import { connect } from 'react-redux';

import { useStoreActions, useStoreState } from 'easy-peasy';

// Externals
import PropTypes from 'prop-types';
import compose from 'recompose/compose';
import validate from 'validate.js';
import _ from 'underscore';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import {
  Grid,
  Button,
  IconButton,
  CircularProgress,
  TextField,
  Typography
} from '@material-ui/core';

// Material icons
import { ArrowBack as ArrowBackIcon } from '@material-ui/icons';

// Shared components
import { Facebook as FacebookIcon, Google as GoogleIcon, AIS } from 'icons';

// Component styles
import styles from './styles';

// Form validation schema
import schema from './schema';

// Store actions
import { loginSuccessful } from "../../store/actions/auth";

const LogIn = (props) => {

    const { classes, history } = props;
    const [values, setValues] = useState({
      email: '',
      password: ''
    });

    const [touched, setTouched] = useState({
      email: false,
      password: false
    });

    const [errors, setErrors] = useState({
      email: null,
      password: null
    });

    const [isValid, setIsValid] = useState(true);
    const isLoading = useStoreState(state => state.auth.isLoading);
    const [submitError, setSubmitError] = useState(null);

    const loginUser = useStoreActions(actions => actions.auth.login);

    const showEmailError = touched.email && errors.email;
    const showPasswordError = touched.password && errors.password;


    const handleBack = () => {
  
      history.goBack();
    };
  
    const validateForm = _.debounce(() => {
      
      const err = validate(values, schema) || {};
  
      setErrors(err);
    }, 300);
  
    const handleFieldChange = (field, value) => {
    
      setSubmitError(null);
      setTouched({
        ...touched,
        [field]: true,
      })
      setValues({
        ...values,
        [field]: value,
      });


    };

    React.useEffect(() => {
      validateForm();
    }, [values]);

    React.useEffect(() => {

      let err = false;
      for (var key in errors) {
        // skip loop if the property is from prototype
        if (!errors.hasOwnProperty(key)) continue;
    
        if(errors[key] !== null) err = true;

      }
      setIsValid(!err)
    }, [errors]);
      
  return (
    <div className={classes.root}>
        <Grid
          className={classes.grid}
          container
        >
          <Grid
            className={classes.quoteWrapper}
            item
            lg={5}
          >
            <div className={classes.quote}>
              <div className={classes.quoteInner}>
                <Typography
                  className={classes.quoteText}
                  variant="h1"
                >
                  ....
                </Typography>
                <div className={classes.person}>
                  <Typography
                    className={classes.name}
                    variant="body1"
                  >
                    ...
                  </Typography>
                  <Typography
                    className={classes.bio}
                    variant="body2"
                  >
                    ...
                  </Typography>
                </div>
              </div>
            </div>
          </Grid>
          <Grid
            className={classes.content}
            item
            lg={7}
            xs={12}
          >
            <div className={classes.content}>
              <div className={classes.contentHeader}>
                <IconButton
                  className={classes.backButton}
                  onClick={handleBack}
                >
                  <ArrowBackIcon />
                </IconButton>
              </div>
              <div className={classes.contentBody}>
                <form className={classes.form}>
                  <Typography
                    className={classes.title}
                    variant="h2"
                  >
                    Prihlásenie
                  </Typography>
                  <Typography
                    className={classes.subtitle}
                    variant="body1"
                  >
                    Prihlásiť sa pomocou treťej strany
                  </Typography>
                  <Button
                    className={classes.googleButton}
                    onClick={() => loginUser(values)}
                    size="large"
                    variant="contained"
                  >
                    <AIS className={classes.googleIcon} />
                    Prihlásenie pomocou AIS konta
                  </Button>
                  <Typography
                    className={classes.sugestion}
                    variant="body1"
                  >
                    alebo sa prihláste pomocou systémového konta
                  </Typography>
                  <div className={classes.fields}>
                    <TextField
                      className={classes.textField}
                      label="Email address"
                      name="email"
                      onChange={event =>
                        handleFieldChange('email', event.target.value)
                      }
                      type="text"
                      value={values.email}
                      variant="outlined"
                    />
                    {showEmailError && (
                      <Typography
                        className={classes.fieldError}
                        variant="body2"
                      >
                        {errors.email[0]}
                      </Typography>
                    )}
                    <TextField
                      className={classes.textField}
                      label="Password"
                      name="password"
                      onChange={event =>
                        handleFieldChange('password', event.target.value)
                      }
                      type="password"
                      value={values.password}
                      variant="outlined"
                    />
                    {showPasswordError && (
                      <Typography
                        className={classes.fieldError}
                        variant="body2"
                      >
                        {errors.password[0]}
                      </Typography>
                    )}
                  </div>
                  {submitError && (
                    <Typography
                      className={classes.submitError}
                      variant="body2"
                    >
                      {submitError}
                    </Typography>
                  )}
                  {isLoading ? (
                    <CircularProgress className={classes.progress} />
                  ) : (
                    <Button
                      className={classes.signInButton}
                      color="primary"
                      disabled={!isValid}
                      onClick={() => loginUser(values)}
                      size="large"
                      variant="contained"
                    >
                      Prilhlásiť sa
                    </Button>
                  )}
                  <Typography
                    className={classes.signUp}
                    variant="body1"
                  >
                    Ešte nemáte konto?{' '}
                    <Link
                      className={classes.signUpUrl}
                      to="/sign-up"
                    >
                      Požiadať
                    </Link>
                  </Typography>
                </form>
              </div>
            </div>
          </Grid>
        </Grid>
      </div>
  );
};
/*
SignIn.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired,
  history: PropTypes.object.isRequired
};*/

export default compose(
  withRouter,
  withStyles(styles)
)(LogIn);
