import React, { Component } from 'react';

import { connect } from 'react-redux';

// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

import { saveRoom, roomSaveInitial } from "store/actions/rooms";

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { SnackbarContent, List, ListItem, ListItemAvatar, Avatar, ListItemText, ListItemSecondaryAction, IconButton } from '@material-ui/core';

import { Delete as DeleteIcon, Folder as FolderIcon, Close as CloseIcon, Warning as WarningIcon } from '@material-ui/icons';
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


const mapStateToProps = state => ({
  roomReducer: state.roomReducer
})

const mapDispatchToProps = dispatch => {
  return {
  }
}


function generate(element) {
  return [0, 1, 2].map(value =>
    React.cloneElement(element, {
      key: value,
    }),
  );
}
class RoomGateways extends Component {

  state = {
  }

  constructor(props) {
    super(props);

  }

  componentDidMount() {
  }
  

  render() {
    const { classes, className, ...rest } = this.props;

    const { savePending, saved, room } = this.props.roomReducer;

    const rootClassName = classNames(classes.root, className);
  
    return (
      <Portlet
        {...rest}
        className={rootClassName}
      >
        <PortletHeader>
          <PortletLabel
            //subtitle="Informácie o miestnosti"
            title="Kontrolné brány"
          />
        </PortletHeader>
        <PortletContent noPadding>
          <div className={classes.demo}>
          {room.gateways.length === 0 &&
          <SnackbarContent
            className={classes.warning}
            aria-describedby="client-snackbar"
            message={
              <span id="client-snackbar" className={classes.message}>
                <WarningIcon className={classes.warningIco} />
                Miestnosť nemá pridelenú žiadnu kontrolnú bránu.
              </span>
            }
          /> }
          {room.gateways.length > 0 &&
            <List >
              {room.gateways.map( g => (
                <ListItem className={classes.gwItem} >
                  <ListItemAvatar>
                    <Avatar>
                      <FolderIcon />
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={(g.name ? g.name : `Nepomenovaná (${g.addr.toUpperCase()})`)}
                    //secondary={secondary ? 'Secondary text' : null}
                  />
                  <ListItemSecondaryAction>
                    <IconButton edge="end" aria-label="delete">
                      <DeleteIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
          </List> }
          </div>
        </PortletContent>
        <PortletFooter className={classes.portletFooter}>
        </PortletFooter>
      </Portlet>
    );
  }
}

RoomGateways.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(RoomGateways);
