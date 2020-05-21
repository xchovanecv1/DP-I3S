import React, { useState } from 'react';

// Externals
import PropTypes from 'prop-types';
import classNames from 'classnames';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, IconButton, CircularProgress } from '@material-ui/core';

// Material icons
import {
  ArrowDownward as ArrowDownwardIcon,
  ArrowUpward as ArrowUpwardIcon,
  Delete as DeleteIcon,
  AddCircleOutlineOutlined as AddIcon
} from '@material-ui/icons';

// Shared components
import { DisplayMode, SearchInput } from 'components';

// Component styles
import styles from './styles';

function Toolbar(props) {

  const [ searchVal, setSearch ] = useState(null);

  const handleSearch = (e) => {
    const val = e.target.value;

    setSearch(val);
  }

  const handleDeleteUsers = () => {

  }


  React.useEffect(() => {
    if (typeof props.search === "function") {
      props.search(searchVal);
    }
  }, [searchVal]);


  const { classes, className, selected, add, onImport } = props;

  const rootClassName = classNames(classes.root, className);

  return (
    <div className={rootClassName}>
      <div className={classes.row}>
        <span className={classes.spacer} />
        {selected.length > 0 && (
          <IconButton
            className={classes.deleteButton}
            onClick={handleDeleteUsers}
          >
            <DeleteIcon />
          </IconButton>
        )}
        <Button
          className={classes.importButton}
          size="small"
          variant="outlined"
          onClick={onImport}
        >
          <ArrowDownwardIcon className={classes.importIcon} /> Import
          </Button>
        <Button
          className={classes.exportButton}
          size="small"
          variant="outlined"
        >
          <ArrowUpwardIcon className={classes.exportIcon} />
          Export
          </Button>
        <Button
          color="primary"
          size="small"
          variant="outlined"
          onClick={() => {
            add();
          }}
        >
          <AddIcon className={classes.addIcon} />
          Pridať
          </Button>
      </div>
      <div className={classes.row}>
        <SearchInput
          className={classes.searchInput}
          placeholder="Vyhľadať zariadenie"
          value={searchVal}
          onChange={handleSearch}
        />
        <CircularProgress className={props.loading && searchVal ? "" : classes.hidden }/>
        <span className={classes.spacer} />
        <DisplayMode mode="list" />
      </div>
    </div>
  );
}

Toolbar.defaultProps = {
  onImport: () => {},
}


export default withStyles(styles)(Toolbar);
