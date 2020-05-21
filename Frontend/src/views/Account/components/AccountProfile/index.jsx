import React, { Component } from 'react';

// Externals
import PropTypes from 'prop-types';
import classNames from 'classnames';

// Material helpers
import { withStyles } from '@material-ui/core';

import { getInitials } from 'helpers';
// Material components
import { Avatar, Typography, Button, LinearProgress } from '@material-ui/core';

// Shared components
import { Portlet, PortletContent, PortletFooter } from 'components';

// Component styles
import styles from './styles';

const AccountProfile = (props) => {

  const { classes, className, token, ...rest } = props;

  const rootClassName = classNames(classes.root, className);
  
  return (
    <Portlet
      {...rest}
      className={rootClassName}
    >
      <PortletContent>
        <div className={classes.details}>
          <div className={classes.info}>
            <Typography variant="h2">{token.sub}</Typography>
            <Typography
              className={classes.locationText}
              variant="body1"
            >
              
            </Typography>
            <Typography
              className={classes.dateText}
              variant="body1"
            >
              
            </Typography>
          </div>
          <Avatar
            className={classes.avatar}
            src="/images/avatars/avatar_11.png"
          >
            {getInitials(token.sub)}
            </Avatar>
        </div>
        <div className={classes.progressWrapper}>
          {/*<Typography variant="body1">Profile Completeness: 70%</Typography>
          <LinearProgress
            value={70}
            variant="determinate"
          />*/}
        </div>
      </PortletContent>
      <PortletFooter>
        <Button
          className={classes.uploadButton}
          color="primary"
          variant="text"
        >
          
        </Button>
        <Button variant="text"></Button>
      </PortletFooter>
    </Portlet>
  );
};

AccountProfile.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(AccountProfile);
