import React, { useState } from 'react';


// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

import {
	KeyboardDatePicker,
} from '@material-ui/pickers';


// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Button, TextField, CircularProgress, Grid } from '@material-ui/core';

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

const PropDetail = (props) => {
	const { classes, className, data, save, newID, add, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [localData, setLocalData] = useState(null);

	const handleInputChange = (event) => {
		const target = event.target;
		const value = target.type === 'checkbox' ? target.checked : target.value;
		const name = target.name;

		setLocalData({
			...localData,
			[name]: value,
		});
	}

	const handleDateChange = date => {
		setLocalData({
			...localData,
			acquired_at: date,
		});
	};

	const saveForm = () => {
		save(localData);
	}

	const addItem = () => {
		add(localData);
	}

	useEffect(() => {
		if (data) {
			setLocalData(data);
		}
	}, [data]);

	if (localData === null) return (<><CircularProgress size={20} color="inherit" /></>);

	const { name, code, price, acquired_at } = localData;
	return (
		<Portlet
			{...rest}
			className={rootClassName}
		>
			<PortletHeader>
				<PortletLabel
					title="Informácie o zariadení"
				/>
			</PortletHeader>
			<PortletContent noPadding>
				<form
					autoComplete="off"
					noValidate
				>
					<Grid
						container
						spacing={3}
					>
						<Grid
							item
							md={6}
							xs={12}
						>
							<div className={classes.field}>
								<TextField
									className={classes.textField}
									name="name"
									label="Názov zariadenia"
									margin="dense"
									required
									fullWidth
									value={name}
									onChange={handleInputChange}
								/>
							</div>
						</Grid>
						<Grid
							item
							md={6}
							xs={12}
						>
							<div className={classes.field}>
								<TextField
									className={classes.textField}
									name="code"
									label="Kódové označenie"
									margin="dense"
									required
									value={code}
									onChange={handleInputChange}
								/>
							</div>
						</Grid>
						<Grid
							item
							md={6}
							xs={12}
						>
							<div className={classes.field}>
								<TextField
									className={classes.textField}
									name="price"
									label="Zostávajúca cena"
									margin="dense"
									value={price}
									type="number"
									onChange={handleInputChange}
								/>
							</div>
						</Grid>
						<Grid
							item
							md={6}
							xs={12}
						>
							<div className={classes.field}>
								<KeyboardDatePicker
									className={classes.textField}
									disableToolbar
									format="dd.MM.yyyy"
									margin="dense"
									id="date-picker"
									label="Dátum zaobstarania"
									value={acquired_at}
									onChange={handleDateChange}
									KeyboardButtonProps={{
										'aria-label': 'upraviť dátum',
									}}
								/>
							</div>
						</Grid>
					</Grid>
				</form>
			</PortletContent>
			<PortletFooter className={classes.portletFooter}>
				{!newID &&
				<Button
					color="primary"
					variant="contained"
					onClick={saveForm}
					className={classes.button}
				//disabled={savePending}
				>
					Uložiť
					{false &&
						<CircularProgress size={20} color="inherit" />
					}
				</Button>
				}
				{newID &&
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
				}
			</PortletFooter>
		</Portlet>
	);
}


PropDetail.propTypes = {
	className: PropTypes.string,
	classes: PropTypes.object.isRequired
};

PropDetail.defaultProps = {
	save: () => { },
}

export default withStyles(styles)(PropDetail);
