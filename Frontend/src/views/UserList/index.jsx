import React, { Component } from 'react';

import { connect } from 'react-redux';
import { loadAllAccounts } from "../../store/actions/account";

// Externals
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { CircularProgress, Typography } from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';

// Shared services
import { getUsers } from 'services/user';

// Custom components
import { UsersToolbar, UsersTable } from './components';

// Component styles
import styles from './style';

const mapStateToProps = state => ({
  ...state
})

const mapDispatchToProps = dispatch => {
  return {
    loadAll: (page, perPage) => dispatch(loadAllAccounts(page + 1, perPage)),
  }
}


class UserList extends Component {
  signal = true;

  state = {
    limit: 10,
    page: 0,
    users: [],
    selectedUsers: [],
    error: null
  };

  async getUsers() {
    this.props.loadAll(this.state.page, this.state.limit);
  }

  componentDidMount() {
    this.signal = true;
    this.getUsers();
  }

  componentWillUnmount() {
    this.signal = false;
  }

  handleSelect = selectedUsers => {
    this.setState({ selectedUsers });
  };

  handleLimitChange = limit => {
    this.setState({ limit }, () => {
      this.getUsers();
    });
  }; 
  handlePageChange = page => {
    this.setState({ page }, () => {
      this.getUsers();
    });
  };

  renderUsers() {
    const { classes } = this.props;
    const { page, limit } = this.state;
    const { accountsPending, accounts, error } = this.props.accountReducer;

    if (accountsPending) {
      return (
        <div className={classes.progressWrapper}>
          <CircularProgress />
        </div>
      );
    }

    if (error) {
      return <Typography variant="h6">{error}</Typography>;
    }

    if (accounts === null || accounts.count === 0) {
      return <Typography variant="h6">There are no users</Typography>;
    }

    return (
      <UsersTable
        //
        onSelect={this.handleSelect}
        onLimitChange={this.handleLimitChange}
        onPageChange={this.handlePageChange}
        users={accounts.accounts}
        count={accounts.count}
        rowsPerPage={limit}
        page={page}
      />
    );
  }

  render() {
    const { classes } = this.props;
    const { selectedUsers } = this.state;
    return (
      <DashboardLayout title="Users">
        <div className={classes.root}>
          <UsersToolbar selectedUsers={selectedUsers} />
          <div className={classes.content}>{this.renderUsers()}</div>
        </div>
      </DashboardLayout>
    );
  }
}

UserList.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(UserList);
