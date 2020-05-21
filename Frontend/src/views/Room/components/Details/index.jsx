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

const Detail = (props) => {
	const { classes, className, data, save, isNew, createItem, ...rest } = props;

	const rootClassName = classNames(classes.root, className);

	const [localData, setLocalData] = useState({
		name: "",
		code: ""
	});

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

	const addNew = () => {
		createItem(localData);
	}

	useEffect(() => {
		if (data) {
			setLocalData(data);
		}
	}, [data]);

	if (localData === null) return (<><CircularProgress size={20} color="inherit" /></>);

	const { name, code } = localData;
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
									label="Názov"
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
									disabled={!isNew}
									value={code}
									onChange={handleInputChange}
								/>
							</div>
						</Grid>
					</Grid>
				</form>
			</PortletContent>
			<PortletFooter className={classes.portletFooter}>
				{!isNew &&
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
				</Button>}
				{isNew &&
				<Button
					color="primary"
					variant="contained"
					onClick={addNew}
					className={classes.button}
				//disabled={savePending}
				>
					Pridať
					{false &&
						<CircularProgress size={20} color="inherit" />
					}
				</Button>}
			</PortletFooter>
		</Portlet>
	);
}


Detail.propTypes = {
	className: PropTypes.string,
	classes: PropTypes.object.isRequired
};

Detail.defaultProps = {
	save: () => { },
}

export default withStyles(styles)(Detail);
