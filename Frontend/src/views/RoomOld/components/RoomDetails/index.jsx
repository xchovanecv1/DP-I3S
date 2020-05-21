import React, { Component } from 'react';

import { connect } from 'react-redux';

// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

//import { saveRoom, roomSaveInitial } from "store/actions/rooms";

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, TextField, CircularProgress } from '@material-ui/core';

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
    //saveRoom: (room, callback) => dispatch(saveRoom(room, callback)),
    //resetRoom: () => dispatch(roomSaveInitial()),
  }
}

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

class Room extends Component {

  state = {
    room: {
      name: "",
    }
  }

  constructor(props) {
    super(props);

    this.handleSave = this.handleSave.bind(this);
    this.roomChanged = this.roomChanged.bind(this);
  }

  roomChanged(room) {
    this.setState({
      room,
    }, () => {
      this.props.history.replace(`/miestnosti/${room.id}`);
    });
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState((state) => {
      return {
        room: {
          ...state.room,
          [name]: value
        },
      }
    });
  }

  componentDidMount() {
    this.setState({
      room: this.props.roomReducer.room,
    });
  }

  handleSave() { 
    this.props.saveRoom(this.state.room, this.roomChanged);
  }


  render() {
    const { classes, className, room, ...rest } = this.props;
    let { name } = this.state.room;

    const { savePending, saved } = this.props.roomReducer;

    const rootClassName = classNames(classes.root, className);
  
    return (
      <Portlet
        {...rest}
        className={rootClassName}
      >
        <PortletHeader>
          <PortletLabel
            subtitle="Informácie o miestnosti"
            title="Miestnosť"
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
                name="name"
                label="Názov miestnosti"
                margin="dense"
                required
                value={name}
                variant="outlined"
                onChange={this.handleInputChange.bind(this)}
              />
            </div>
          </form>
        </PortletContent>
        <PortletFooter className={classes.portletFooter}>
          <Button
            color={(!saved ? "primary" : "secondary")}
            variant="contained"
            onClick = {this.handleSave}
            className={classes.button}
            disabled={savePending}
          >
            {!savePending && !saved && `Uložiť`}
            {!savePending && saved && `Uložené`}
            {savePending &&
                <CircularProgress size={20} color="inherit" />
            }
          </Button>
        </PortletFooter>
      </Portlet>
    );
  }
}

Room.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Room);
