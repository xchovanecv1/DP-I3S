import React, { Component, Fragment, useState } from 'react';

// Externals
import classNames from 'classnames';
import compose from 'recompose/compose';
import PropTypes from 'prop-types';

// Material helpers
import { withStyles, withWidth } from '@material-ui/core';

// Material components
import { Drawer, Snackbar, SnackbarContent} from '@material-ui/core';

// Custom components
import { Sidebar, Topbar, Footer } from './components';

import PropDialog from '../../components/PropertyPosses';

// Component styles
import styles from './styles';
import { useEffect } from 'react';


const Dashboard = (props) => {
  
  
  const [isOpen, setIsOpen] = useState(true);
  const [inited, setInited] = useState(false);

  const [possesOpen, setPossesOpen] = useState(false);

  const isMobile = ['xs', 'sm', 'md'].includes(props.width);

  const pO = localStorage.getItem("posses-ignore")

  useEffect(() => {
    if (!inited) {
      setIsOpen(!isMobile);
      setInited(true);

      if (!pO) {
        setPossesOpen(true);
      }
    }
  }, []);


  const shiftTopbar = isOpen && !isMobile;
  const shiftContent = isOpen && !isMobile;

  const handleClose = () => {
    setIsOpen(false);
  };

  const handleToggleOpen = () => {
    setIsOpen(!isOpen);
  };

  const closePosses = () => {
    localStorage.setItem("posses-ignore", true)
    setPossesOpen(false);
  }

  const clearIgnore = () => {
    localStorage.removeItem("posses-ignore")
    setPossesOpen(true);
  }

  const { classes, width, title, children } = props;

  return (
    <Fragment>
      <Topbar
        className={classNames(classes.topbar, {
          [classes.topbarShift]: shiftTopbar
        })}
        isSidebarOpen={isOpen}
        onToggleSidebar={handleToggleOpen}
        title={title}
      />
      <Drawer
        anchor="left"
        classes={{ paper: classes.drawerPaper }}
        onClose={handleClose}
        open={isOpen}
        variant={isMobile ? 'temporary' : 'persistent'}
      >
        <Sidebar className={classes.sidebar} />
      </Drawer>
      <main
        className={classNames(classes.content, {
          [classes.contentShift]: shiftContent
        })}
      >
        {children}
        <PropDialog 
          open={possesOpen}
          close={closePosses}
          clearIgnore={clearIgnore}
        />
        <Footer />
      </main>
    </Fragment>
  );

}

Dashboard.propTypes = {
  children: PropTypes.node,
  className: PropTypes.string,
  classes: PropTypes.object.isRequired,
  title: PropTypes.string,
  width: PropTypes.string.isRequired
};

export default compose(
  withStyles(styles),
  withWidth()
)(Dashboard);
