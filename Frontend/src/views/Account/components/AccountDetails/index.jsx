import React, { Component } from 'react';

// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, TextField } from '@material-ui/core';

// Shared components
import {
  Portlet,
  PortletHeader,
  PortletLabel,
  PortletContent,
  PortletFooter
} from 'components';

// Component styles
import styles from './styles';

const states = [
  {
    value: 'alabama',
    label: 'Alabama'
  },
  {
    value: 'new-york',
    label: 'New York'
  },
  {
    value: 'san-francisco',
    label: 'San Francisco'
  }
];

const Account = (props) => {
  const { classes, className, token, ...rest } = props;
  const { phone, state, country, email } = props;

  const rootClassName = classNames(classes.root, className);

  const handleChange = e => {
    this.setState({
      state: e.target.value
    });
  };

  const parsed = token.sub && token.sub.split(" ");
  const firstName = parsed[0];
  const lastName = parsed[1];

  return (
    <Portlet
      {...rest}
      className={rootClassName}
    >
      <PortletHeader>
        <PortletLabel
          subtitle="Základné informácie o Vás"
          title="Profil"
        />
      </PortletHeader>
      <PortletContent noPadding>
        <form
          autoComplete="off"
          noValidate
        >
          <div className={classes.field}>
            <TextField
              className={classes.textField}
               label="Meno"
              margin="dense"
              required
              value={firstName}
              variant="outlined"
            />
            <TextField
              className={classes.textField}
              label="Priezvisko"
              margin="dense"
              required
              value={lastName}
              variant="outlined"
            />
          </div>
        </form>
      </PortletContent>
      <PortletFooter className={classes.portletFooter}>
        <Button
          color="primary"
          variant="contained"
        >
          Uložiť
          </Button>
      </PortletFooter>
    </Portlet>
  );

}

Account.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Account);
