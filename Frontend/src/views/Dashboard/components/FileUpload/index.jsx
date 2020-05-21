import React, { Component } from 'react';

import { useStoreActions } from 'easy-peasy';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Typography, LinearProgress } from '@material-ui/core';

// Material icons
import { InsertChartOutlined as InsertChartIcon } from '@material-ui/icons';

// Shared components
import { Paper } from 'components';

// Component styles
import styles from './styles';

const FileUpload = (props) => {
  const { classes, className, ...rest } = props;

  const rootClassName = classNames(classes.root, className);
  const upload = useStoreActions(actions => actions.import.upload);
  return (<Paper
    {...rest}
    className={rootClassName}
  >
    <div className={classes.content}>
      <div className={classes.details}>
        <form
              enctype="multipart/form-data"
              action="/api/admin/upload"
              method="post"
            >
              <input type="file" name="myFile" />
              <input type="submit" value="upload" />
            </form>
            <input type="file" name="file" onChange={upload}/>
      </div>
      <div className={classes.iconWrapper}>
        <InsertChartIcon className={classes.icon} />
      </div>
    </div>
    <div className={classes.footer}>
      <LinearProgress
        value={75.5}
        variant="determinate"
      />
    </div>
  </Paper>);
}
export default withStyles(styles)(FileUpload);
