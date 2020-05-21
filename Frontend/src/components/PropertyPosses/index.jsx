import React, { useState } from 'react';

import { useStoreActions, useStoreState } from "easy-peasy";
//import Autocomplete from '@material-ui/lab/Autocomplete';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';
import moment from 'moment';

import Autocomplete from '@material-ui/lab/Autocomplete';

import { Warning as WarningIcon, SimCard, MemoryOutlined, Delete } from "@material-ui/icons";

import {
	KeyboardDatePicker,
} from '@material-ui/pickers';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, Table, TableBody, TableHead, TableRow, TableCell, TextField, Snackbar, CircularProgress, Link, Typography, IconButton, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@material-ui/core';

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

const PropPosses = (props) => {
	const { classes, className, id, save, open, close, clearIgnore, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	//confirm

	const loadPending = useStoreActions(actions => actions.props.loadPendinProps);
	const pendingProps = useStoreState(state => state.props.possesPending);

	const savePending = useStoreActions(actions => actions.props.savePendinProps);

	const [inited, setInited] = useState(false);

	const [localPending, setLocalPending] = useState(null);

	useEffect(() => {
		if (!inited) {
			loadPending();
			setInited(true);
		}
	}, []);

	useEffect(() => {
		if (Array.isArray(pendingProps) && pendingProps.length) {
			setLocalPending(pendingProps);
			setInited(true);
		} else {
			setLocalPending(null);
		}
	}, [pendingProps]);

	const changeValue = (event, i) => {
		const target = event.target;
		const value = target.type === 'checkbox' ? target.checked : target.value;

		const bf = localPending[i];
		bf.comment = value;

		localPending[i] = bf;

		setLocalPending([
			...localPending,
		]);
	}

	const confirm = () => {
		savePending(localPending).then(() => {
			loadPending();
		});
	}

	const isPending = Array.isArray(localPending) && localPending.length;
	console.log(localPending);
	return (
		<>
			<Dialog
				open={open && isPending}
				onClose={close}
				aria-labelledby="alert-dialog-title"
				aria-describedby="alert-dialog-description"
			>
				<DialogTitle id="alert-dialog-title">Potvrdenie vypožičania zariadení</DialogTitle>
				<DialogContent>
					<Table>
						<TableHead>
							<TableRow>
								<TableCell align="left">
									Zariadenie
								</TableCell>
								<TableCell align="left">
									Poznámka
								</TableCell>
							</TableRow>
						</TableHead>
						<TableBody>
							{
								Array.isArray(localPending) && localPending.length && localPending.map((p, i) => (
									<TableRow>
										<TableCell align="left">
											{p.property_name} ({p.property_code})
										</TableCell>
										<TableCell align="left">
											<TextField
												value={p.comment}
												onChange={(e) => { changeValue(e, i)}}
											/>
										</TableCell>
									</TableRow>
								))
							}
						</TableBody>
					</Table>
				</DialogContent>
				<DialogActions>
					<Button onClick={close} color="primary">
						Neskôr
				</Button>
					<Button onClick={() => {confirm()}} color="primary">
						Potvrdiť
				</Button>
				</DialogActions>
			</Dialog>

			<Snackbar
			anchorOrigin={{ vertical:'bottom', horizontal: "right" }}
			key={`bottom,right`}
			open={!open && isPending }
			onClick={() => { clearIgnore(); console.log("test") }}
			autoHideDuration={6000}
			ContentProps={{
				'aria-describedby': 'message-id',
			}}
			message={<span id="message-id">Máte nepotvrdené vypožičania</span>}
			/>
		</>
	);
}


PropPosses.propTypes = {
	className: PropTypes.string,
	classes: PropTypes.object.isRequired
};

PropPosses.defaultProps = {
	save: () => { },
};

export default withStyles(styles)(PropPosses);
