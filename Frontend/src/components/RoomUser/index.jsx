import React, { useState } from 'react';

import { useStoreActions, useStoreState } from "easy-peasy";
//import Autocomplete from '@material-ui/lab/Autocomplete';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

import { Warning as WarningIcon, InfoOutlined as InfoIcon, CheckCircleOutline, HowToRegOutlined } from "@material-ui/icons";

import {
	KeyboardDatePicker,
} from '@material-ui/pickers';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, TextField, CircularProgress, Link, Typography, IconButton, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@material-ui/core';

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
import { useEffect } from 'react';

const Picker = (props) => {
	const { classes, className, room, roomPending, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [inited, setInited] = useState(false);

	const hasUser = room && room.user;
	const user = room && room.user;

	useEffect(() => {
		if (!inited) {
			setInited(true);
		}
	}, []);

	
	return (
				<>
					<Portlet
						{...rest}
						className={rootClassName}
					>
						<PortletHeader>
							<PortletLabel
								title="Aktívny používateľ"
							/>
						</PortletHeader>
						<PortletContent noPadding>
							{
								roomPending && 
								<><CircularProgress size={20} color="inherit" /></>
							}
							{
								!roomPending && !hasUser &&
								<div className={classes.field}>
									<span id="client-snackbar" className={classes.message}>
										<InfoIcon className={classes.warningIco} />
								V miestnosti nie je žiadny autorizovaný používateľ
							</span>
								</div>
							}
							{
								!roomPending && hasUser &&
								<div className={classes.field}>
									<span id="client-snackbar" className={classes.message}>
										<InfoIcon className={classes.userIco} />
										{`V miestnosti je aktuálne autorizovaný používateľ`}
										{
											(
												<>
												&nbsp;
												<Link to={`/pouzivatelia/${user.id}`}>
														<Typography
															className={classes.nameText}
															variant="body1"
														>
															{user.name}
														</Typography>
												</Link>
												</>
											)
										}
									</span>
								</div>
							}
						</PortletContent>
					</Portlet>
				</>
	);
}


Picker.propTypes = {
					className: PropTypes.string,
	classes: PropTypes.object.isRequired
};

Picker.defaultProps = {
					save: () => {},
};

export default withStyles(styles)(Picker);
