import React, { useState } from 'react';

import { useStoreActions, useStoreState } from "easy-peasy";
//import Autocomplete from '@material-ui/lab/Autocomplete';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';
import moment from 'moment';

import Autocomplete from '@material-ui/lab/Autocomplete';

import { Warning as WarningIcon, SimCard, MemoryOutlined, Delete, RssFeedOutlined} from "@material-ui/icons";

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
	const [editing, setEditing] = useState(false);

	const items = useStoreState(state => state.rooms.gateways);
	const itemsError = useStoreState(state => state.rooms.gatewaysErrpr);
	const itemsLoading = useStoreState(state => state.rooms.gatewaysLoading);


	const [itemUnlink, setItemUnlink] = useState(null);

	const [open, setOpen] = React.useState(false);

	const loadItems = useStoreActions(actions => actions.rooms.loadGateways);

	const performSearch = useStoreActions(actions => actions.rooms.searchGateways);
	const searchRes = useStoreState(state => state.rooms.gws);
	const loading = useStoreState(state => state.cards.gwSearchLoading);

	const performLinkDelete = useStoreActions(actions => actions.rooms.deleteGateway);

	const itemAdding = useStoreActions(actions => actions.rooms.gwAddPending);
	const itemAdd = useStoreActions(actions => actions.rooms.addRoomGw);
/*
	const fromLogAdding = useStoreState(state => state.cards.cardLinkAdding);
*/	
	const [ picked, setPicked ] = useState(null);

	const options = searchRes && searchRes.gateways || [];

	const handleSearch = (event) => {
		const value = event.target && event.target.value;
		if (value) {
			performSearch(value);
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

	const addItem = () => {
		console.log('card adding', picked, id);

		itemAdd({
			id,
			gateway: picked
		}).then(() => {
			loadItems(id);
			setEditing(false);
		});

		/*
		addFromLog({
			prop_id: id,
			data: picked,
		}).then((r) => {
			loadCards(id);
			setEditing(false);
		});*/
	}
	const deleteLink = (iid) => {
		setItemUnlink(iid);
		console.log("deletin", iid);
	}

	const clearLink = () => {
		setItemUnlink(null);
	}
	const performUnlink = () => {
		//setCardUnlink(null);
		const data = {
			prop_id: id,
			data: itemUnlink,
		};
		performLinkDelete(data).then((r) => {
			console.log("deleted ", r);
			loadItems(id);
			setItemUnlink(null);
		});
	}

	const renderItems = () => {
	
		return (
		  <>
			{items.map((card, i) => (
			  <div
				className={classes.product}
				key={i}
			  >
				<div className={classes.productImageWrapper}>
				  <RssFeedOutlined
					alt="Product Name"
					className={classes.productImage}
					//src={product.imageUrl}
				  />
				</div>
				<div className={classes.productDetails}>
				  <Link to="#">
					<Typography
					  className={classes.productTitle}
					  variant="h5"
					>
					  {card.addr}
					</Typography>
				  </Link>
				  <Typography
					className={classes.productTimestamp}
					variant="body2"
				  >
					{card.name}
				  </Typography>
				</div>
				<div>
				  <IconButton>
					<Delete 
						onClick={() => deleteLink(card)}
					/>
				  </IconButton>
				</div>
			  </div>
			))}
		  </>
		);
	  }


	useEffect(() => {
		if (id && !inited) {
			setInited(true);
			loadItems(id);
		}
	}, [id]);

	useEffect(() => {
		if (editing) {
			performSearch("");
		}
	}, editing);

	if (itemsLoading) return (<><CircularProgress size={20} color="inherit" /></>);

	const itemHasData = items && items.length > 0;

	return (
		<>
			<Portlet
				{...rest}
				className={rootClassName}
			>
				<PortletHeader>
					<PortletLabel
						title="Informácie o prístupových bodoch"
					/>
				</PortletHeader>
				<PortletContent noPadding>
					{itemHasData &&
						renderItems()
					}
					{
						!itemHasData && !itemsLoading &&
						<div className={classes.field}>
							<span id="client-snackbar" className={classes.message}>
								<WarningIcon className={classes.warningIco} />
								Miestnosť nemá pridelené kontrolné miesto
							</span>
						</div>
					}
					{editing &&
						<div className={classes.field}>
							{itemAdding && <><CircularProgress size={20} color="inherit" /></>}
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
								disabled={itemAdding}
								onChange={handleSearchPick} 
								getOptionSelected={(option, value) => option.code === value.code}
								getOptionLabel={(option) => `${option.addr} (${moment(option.last_active).format('HH:mm:ss')})`}
								noOptionsText="Žiadny záznam"
								loadingText="Vyhľadávanie..."
								options={options}
								loading={loading}
								renderInput={(params) => (
									<TextField
										{...params}
										label="Prístupová brána"
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
							onClick={() => { setEditing(true); performSearch(""); }}
							className={classes.button}
						//disabled={savePending}
						>
							Pridať
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
								onClick={addItem}
								className={classes.button}
							//disabled={savePending}
							>
								Pridať
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
			<Dialog
				open={itemUnlink}
				//onClose={handleClose}
				aria-labelledby="alert-dialog-title"
				aria-describedby="alert-dialog-description"
			>
				<DialogTitle id="alert-dialog-title">Odstránenie identifikátora zariadenia</DialogTitle>
				<DialogContent>
				<DialogContentText id="alert-dialog-description">
					Chystáte sa odstrániť identifikátor priradený k tomuto zariadeniu, čím znemožníte jeho ďaľšie sledovanie. Prajete si pokračovať?
				</DialogContentText>
				</DialogContent>
				<DialogActions>
				<Button onClick={clearLink} color="primary">
					Nie
				</Button>
				<Button onClick={performUnlink} color="primary" autoFocus>
					Odstrániť
				</Button>
				</DialogActions>
			</Dialog>
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
