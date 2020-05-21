import React, { useState } from 'react';

// Externals
import PropTypes from 'prop-types';
import classNames from 'classnames';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, IconButton } from '@material-ui/core';

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


  const { classes, className, confirm, close } = props;

  const rootClassName = classNames(classes.root, className);

  return (
    <div className={rootClassName}>
      <div className={classes.row}>
        <span className={classes.spacer} />
        <Button
          className={classes.exportButton}
          size="small"
          variant="outlined"
          onClick={() => {
            close();
          }}
        >
          <DeleteIcon />
          Zrušiť
          </Button>
        <Button
          color="primary"
          size="small"
          variant="outlined"
          onClick={() => {
            confirm();
          }}
        >
          <AddIcon className={classes.addIcon} />
          Importovať
          </Button>
      </div>
    </div>
  );
}

Toolbar.defaultProps = {
  confirm: () => {},
  close: () => {},
}


export default withStyles(styles)(Toolbar);
