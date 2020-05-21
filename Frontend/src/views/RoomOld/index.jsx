import React, { Component } from 'react';

import moment from 'moment';
import { connect } from 'react-redux';
//import { loadRoom } from "../../store/actions/rooms";
// Externals
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Grid, CircularProgress } from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';

// Custom components
import RoomDetails from './components/RoomDetails';
import RoomGateways from './components/RoomGateways';

// Component styles
const styles = theme => ({
  root: {
    padding: theme.spacing.unit * 4
  }
});

const mapStateToProps = state => ({
  roomReducer: state.roomReducer
})

const mapDispatchToProps = dispatch => {
  return {
    //loadRoom: (id) => dispatch(loadRoom(id)),
  }
}

class Room extends Component {
  state = { tabIndex: 0 };

  componentDidMount() {
    const { params } = this.props.match;
    
    if (params.id) {
      this.props.loadRoom(params.id);
    }

  }

  render() {
    const { classes, history } = this.props;
    const { roomReducer } = this.props;
    const { roomPending, room, saved } = roomReducer;
    const { params } = this.props.match;
    
    return (
      <DashboardLayout title="Room">
        <div className={classes.root}>
          {roomPending && 
            <div className={classes.progressWrapper}>
              <CircularProgress />
            </div>
          }
          {!roomPending &&
          <>
            <Grid
              container
              spacing={4}
            >
              <Grid
                item
                lg={8}
                md={6}
                xl={8}
                xs={12}
              >
                <RoomDetails history={history}/>
              </Grid>
              <Grid
                item
                lg={4}
                md={6}
                xl={4}
                xs={12}
              >
                <RoomGateways />
              </Grid>
            </Grid>
            <Grid
            container
            spacing={4}
          >
            <Grid
              item
              lg={8}
              md={6} 
              xl={8}
              xs={12}
            >
              <RoomDetails history={history}/>
            </Grid>
            <Grid
              item
              lg={4}
              md={6}
              xl={4}
              xs={12}
            >
            </Grid>
          </Grid>
          </>
        }
        </div>
      </DashboardLayout>
    );
  }
}

Room.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Room);
