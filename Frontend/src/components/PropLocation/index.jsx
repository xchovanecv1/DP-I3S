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
	const { classes, className, id, save, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [inited, setInited] = useState(false);

	const loadLocation = useStoreActions(actions => actions.props.loadLocation)
	const locationLoading = useStoreState(state => state.props.locationLoading)
	const location = useStoreState(state => state.props.location);

	const hasLocation = location && location.type;

	useEffect(() => {
		if (!inited) {
			loadLocation(id);
			setInited(true);
		}
	}, []);

	if (locationLoading) {
		return <CircularProgress size={20} color="inherit" />;
	}

	console.log(location);

	const renderExit = () => {

		const user = location.user;
		const room = location.room;
		const lastUser = user && user.name;
		const lastRoom = room && room.name;

		return (
			<span id="client-snackbar" className={classes.message}>
				<WarningIcon className={classes.warningIco} />
				{`Zariadenie sa nenachádza v žiadnej miestnosti.`}
				{
					lastUser && (
						<>
							Prevzaté používateľom&nbsp;
							<Link to={`/pouzivatelia/${user.id}`}>
								<Typography
									className={classes.nameText}
									variant="body1"
								>
									{user.name}
								</Typography>
							</Link>
							&nbsp;
						</>
					)
				}
				{
					!lastUser && (
						<>
							Prevzaté neautorizovaným používateľom&nbsp;
						</>
					)
				}
				{
					lastRoom && (
						<>
							z miestnosti&nbsp;
							<Link to={`/miesnosti/${room.id}`}>
								<Typography
									className={classes.nameText}
									variant="body1"
								>
									{room.name} ({room.code})
								</Typography>
							</Link>
						</>
					)
				}
			</span>
					);
	}

	const renderPosses = () => {

		const user = location.user;
		//const room = location.room;
		const lastUser = user && user.name;
		//const lastRoom = room && room.name;
		const reason = location && location.reason;

		return (
			<span id="client-snackbar" className={classes.message}>
				<HowToRegOutlined className={classes.userIco} />
				{`Zariadenie je aktuálne vypožičané používateľom`}
				{
					lastUser && (
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
				{
					reason && (
						<>
							(Poznámka:&nbsp;<Typography
									className={classes.nameText}
									variant="body1"
								>
								{reason}
							</Typography>)
						</>
					)
				}
			</span>
					);
	}

	const renderEnter = () => {

		//const user = location.user;
		const room = location.room;
		//const lastUser = user && user.name;
		const lastRoom = room && room.name;

		return (
			<span id="client-snackbar" className={classes.message}>
				<CheckCircleOutline className={classes.enterIco} />
				{`Zariadenie sa nachádza v miestnosti`}
				{
					lastRoom && (
						<>
							&nbsp;
							<Link to={`/miesnosti/${room.id}`}>
								<Typography
									className={classes.nameText}
									variant="body1"
								>
									{room.name} ({room.code})
								</Typography>
							</Link>
						</>
					)
				}
			</span>
					);
	}

	const renderLocation = () => {
		return (
				<div className={classes.field}>
					{
						location.type === "EXIT" && renderExit()
					}
					{
						location.type === "ENTER" && renderEnter()
					}
					{
						location.type === "POSSES" && renderPosses()
					}
				</div>
		);
	}

	return (
				<>
					<Portlet
						{...rest}
						className={rootClassName}
					>
						<PortletHeader>
							<PortletLabel
								title="Pozícia zariadenia"
							/>
						</PortletHeader>
						<PortletContent noPadding>

							{
								!hasLocation &&
								<div className={classes.field}>
									<span id="client-snackbar" className={classes.message}>
										<WarningIcon className={classes.warningIco} />
								K zariadeniu nie je akutálne registrovaná žiadna poloha
							</span>
								</div>
							}
							{
								hasLocation && renderLocation()
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
