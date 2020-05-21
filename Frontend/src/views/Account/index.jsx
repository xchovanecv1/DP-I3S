import React, { Component } from 'react';

import { useStoreActions } from "easy-peasy";
// Externals
import PropTypes from 'prop-types';

import { useStoreState } from "easy-peasy";
// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Grid } from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';

// Custom components
import { AccountProfile, AccountDetails } from './components';

import CardDetail from '../../components/CardPicker';
import CardPicker from 'components/CardPicker';

// Component styles
const styles = theme => ({
  root: {
    padding: theme.spacing.unit * 4
  }
});

//
const Account = (props) => {

  const { classes } = props;

  const token = useStoreState(state => state.auth.token);

	const performLinkDelete = useStoreActions(actions => actions.accounts.deleteCardLink);

	const addFromLog = useStoreActions(actions => actions.accounts.addCardLinkFromLog);
	const fromLogAdding = useStoreState(state => state.accounts.cardLinkAdding);
	
	const loadCards = useStoreActions(actions => actions.cards.loadCard);

  const id = token.id;

  return (
    <DashboardLayout title="Account">
      <div className={classes.root}>
        <Grid
          container
          spacing={4}
        >
          <Grid
            item
            lg={4}
            md={6}
            xl={4}
            xs={12}
          >
            <AccountProfile 
              token={token}
            />
          </Grid>
          <Grid
            item
            lg={8}
            md={6}
            xl={8}
            xs={12}
          >
            <AccountDetails 
              token={token}/>
          </Grid>
        </Grid>
      </div>
      <div className={classes.root}>
        <Grid
          container
          spacing={4}
        >
          <Grid
            item
            lg={4}
            md={6}
            xl={4}
            xs={12}
          >
            <CardPicker 
              label="Prístupové karty"
              emptyText="Nemáte priradenú žiadnu kartu"
							performLinkDelete={performLinkDelete}
							addFromLog={addFromLog}
							fromLogAdding={fromLogAdding}
							loadCards={loadCards}
              controllerContext="accounts"
              id={id}
            />
          </Grid>
          <Grid
            item
            lg={8}
            md={6}
            xl={8}
            xs={12}
          >
          </Grid>
        </Grid>
      </div>
    </DashboardLayout>
  );
}

Account.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Account);
