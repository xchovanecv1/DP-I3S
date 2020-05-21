import React, { useState } from 'react';

import { useStoreActions, useStoreState } from "easy-peasy";
//import Autocomplete from '@material-ui/lab/Autocomplete';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

import Autocomplete from '@material-ui/lab/Autocomplete';

import { Warning as WarningIcon } from "@material-ui/icons";

import {
	KeyboardDatePicker,
} from '@material-ui/pickers';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, TextField, CircularProgress } from '@material-ui/core';

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

const RoomPicker = (props) => {
	const { classes, className, id, save, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [inited, setInited] = useState(false);
	const [editing, setEditing] = useState(false);

	const room = useStoreState(state => state.rooms.room);
	const roomError = useStoreState(state => state.rooms.roomError);
	const roomLoading = useStoreState(state => state.rooms.roomLoading);

	const [open, setOpen] = React.useState(false);

	const loadRoom = useStoreActions(actions => actions.rooms.loadRoom);
	const clearRoom = useStoreActions(actions => actions.rooms.clearRoom);

	const searchRooms = useStoreActions(actions => actions.rooms.searchRooms);
	const searchRes = useStoreState(state => state.rooms.search);
	const loading = useStoreState(state => state.rooms.searchLoading);
	
	const [ picked, setPicked ] = useState(null);

	const options = searchRes && searchRes.rooms || [];

	const handleSearch = (event) => {
		const value = event.target && event.target.value;
		if (value) {
			searchRooms(value);
		}
	}

	const handleInputChange = (event) => {
		const target = event.target;
		const value = target.type === 'checkbox' ? target.checked : target.value;
		const name = target.name;
		/*
				setLocalData({
					...localData,
					[name]: value,
				});*/
	}

	const handleSearchPick = (evt, val) => {
		setPicked(val);
	}

	const saveRoom = () => {
		save(picked);
		setEditing(false);
	}

	useEffect(() => {
		if (!inited) {
			if(id) {
				loadRoom(id);
				searchRooms("");
			} else {
				clearRoom();
				searchRooms("");
			}
			setInited(true);
		}
	}, []);

	if (roomLoading) return (<><CircularProgress size={20} color="inherit" /></>);

	const roomHasData = room.name !== null && room.name.length > 0;

	const { name, code } = room;
	return (
		<Portlet
			{...rest}
			className={rootClassName}
		>
			<PortletHeader>
				<PortletLabel
					title="Informácie o miestnosti"
				/>
			</PortletHeader>
			<PortletContent noPadding>
				{!editing && roomHasData &&
					<form
						autoComplete="off"
						noValidate
					>
						<div className={classes.field}>
							<TextField
								className={classes.textField}
								name="name"
								label="Názov miestnosti"
								margin="dense"
								value={name}
								inputProps={{
									readOnly: true,
									disabled: true,
								}}
								onChange={(v) => console.log(v)}
							/>
						</div>
						<div className={classes.field}>
							<TextField
								className={classes.textField}
								name="code"
								label="Kód miestnosti"
								margin="dense"
								inputProps={{
									readOnly: true,
									disabled: true,
								}}
								value={code}
							/>
						</div>
					</form>
				}
				{
					!editing && !roomHasData &&
					<div className={classes.field}>
						<span id="client-snackbar" className={classes.message}>
							<WarningIcon className={classes.warningIco} />
							Zariadenie nemá pridelenú žiadnu miestnosť
						</span>
					</div>
				}
				{editing &&
					<div className={classes.field}>
						<Autocomplete
							id="asynchronous-search"
							style={{ width: 300 }}
							open={open}
							onOpen={() => {
								setOpen(true);
							}}
							onClose={() => {
								setOpen(false);
							}}
							onChange={handleSearchPick} 
							getOptionSelected={(option, value) => option.name === value.name}
							getOptionLabel={(option) => `${option.name} (${option.code})`}
							options={options}
							loading={loading}
							renderInput={(params) => (
								<TextField
									{...params}
									label="Názov miestnosti"
									onChange={handleSearch}
									InputProps={{
										...params.InputProps,
										endAdornment: (
											<React.Fragment>
												{loading ? <CircularProgress color="inherit" size={20} /> : null}
												{params.InputProps.endAdornment}
											</React.Fragment>
										),
									}}
								/>
							)}
						/>
					</div>
				}
			</PortletContent>
			<PortletFooter className={classes.portletFooter}>
				{!editing &&
					<Button
						color="primary"
						variant="contained"
						onClick={() => { setEditing(true) }}
						className={classes.button}
					//disabled={savePending}
					>
						Upraviť
						{false &&
							<CircularProgress size={20} color="inherit" />
						}
					</Button>
				}
				{editing &&
					<>
						<Button
							color="primary"
							variant="contained"
							onClick={saveRoom}
							className={classes.button}
						//disabled={savePending}
						>
							Uložiť
							{false &&
								<CircularProgress size={20} color="inherit" />
							}
						</Button>
						<Button
							color="secondary"
							variant="contained"
							onClick={() => {
								setEditing(false);
							}}
							className={classes.button}
						//disabled={savePending}
						>
							Zrušiť
							{false &&
								<CircularProgress size={20} color="inherit" />
							}
						</Button>
					</>
				}
			</PortletFooter>
		</Portlet>
	);
}


RoomPicker.propTypes = {
	className: PropTypes.string,
	classes: PropTypes.object.isRequired
};

RoomPicker.defaultProps = {
	save: () => {},
};

export default withStyles(styles)(RoomPicker);
