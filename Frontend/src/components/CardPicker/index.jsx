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
	const { classes, className, id, save, label, emptyText, performLinkDelete, addFromLog, fromLogAdding, loadCards, controllerContext, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [inited, setInited] = useState(false);
	const [editing, setEditing] = useState(false);

	const cards = useStoreState(state => state.cards.card);
	const cardError = useStoreState(state => state.cards.cardError);
	const cardLoading = useStoreState(state => state.cards.cardLoading);


	const [cardUnlink, setCardUnlink] = useState(null);

	const [open, setOpen] = React.useState(false);

	const performSearch = useStoreActions(actions => actions.cards.searchCardLogs);
	const searchRes = useStoreState(state => state.cards.logSearch);
	const loading = useStoreState(state => state.cards.logSearchLoading);

	const [ picked, setPicked ] = useState(null);

	const options = searchRes && searchRes.cards || [];

	const handleSearch = (event) => {
		const value = event.target && event.target.value;
		if (value) {
			performSearch(value);
		} else {
			performSearch("");
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

	const addCard = () => {
		console.log('card adding', picked);
		addFromLog({
			id: id,
			data: picked,
		}).then((r) => {
			loadCards({
				ctx: controllerContext,
				id:id
			});
			setEditing(false);
		});
	}
	const deleteLink = (cardId) => {
		setCardUnlink(cardId);
		console.log("deletin", cardId);
	}

	const clearLink = () => {
		setCardUnlink(null);
	}
	const performUnlink = () => {
		//setCardUnlink(null);
		const data = {
			id: id,
			data: cardUnlink,
		};
		performLinkDelete(data).then((r) => {
			console.log("deleted ", r);
			loadCards({
				ctx: controllerContext,
				id:id
			});
			setCardUnlink(null);
		});
	}

	const renderCards = () => {
	
		return (
		  <>
			{cards.cards.map((card, i) => (
			  <div
				className={classes.product}
				key={i}
			  >
				<div className={classes.productImageWrapper}>
				  <MemoryOutlined
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
					  {card.code}
					</Typography>
				  </Link>
				  <Typography
					className={classes.productTimestamp}
					variant="body2"
				  >
					{card.type}
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
			loadCards({
				ctx: controllerContext,
				id:id
			});
		}
	}, []);

	useEffect(() => {
		if (editing) {
			performSearch("");
		}
	}, editing);


	if (cardLoading) return (<><CircularProgress size={20} color="inherit" /></>);

	const cardHasData = cards.cards && cards.cards.length > 0;

	const code = "a";
	const type = "s";
	return (
		<>
			<Portlet
				{...rest}
				className={rootClassName}
			>
				<PortletHeader>
					<PortletLabel
						title={label}
					/>
				</PortletHeader>
				<PortletContent noPadding>
					{cardHasData &&
						renderCards()
					}
					{
						!cardHasData &&
						<div className={classes.field}>
							<span id="client-snackbar" className={classes.message}>
								<WarningIcon className={classes.warningIco} />
								{emptyText}
							</span>
						</div>
					}
					{editing &&
						<div className={classes.field}>
							{fromLogAdding && <><CircularProgress size={20} color="inherit" /></>}
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
								disabled={fromLogAdding}
								onChange={handleSearchPick} 
								getOptionSelected={(option, value) => option.code === value.code}
								getOptionLabel={(option) => `${option.code} (${moment(option.created_at).format('HH:mm:ss')})`}
								noOptionsText="Žiadny záznam"
								loadingText="Vyhľadávanie..."
								options={options}
								loading={loading}
								renderInput={(params) => (
									<TextField
										{...params}
										label="Číslo nepoužitej karty"
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
								onClick={addCard}
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
				open={cardUnlink}
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
