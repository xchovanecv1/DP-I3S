import React, { useState } from 'react';
import { Link } from 'react-router-dom';

// Externals
import classNames from 'classnames';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';

// Material helpers
import { withStyles } from '@material-ui/core';

// Dialog helpers
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, FormControlLabel } from '@material-ui/core';

// Material components
import {
	Avatar,
	Checkbox,
	Table,
	TableBody,
	CircularProgress,
	TableCell,
	TableHead,
	TableRow,
	FormControl,
	MenuItem,
	InputLabel,
	Select,
	FormHelperText,
	Typography,
	TablePagination,
	Button
} from '@material-ui/core';

// Shared helpers
import { getInitials } from 'helpers';

// Shared components
import { Portlet, PortletContent } from 'components';

// Component styles
import styles from './styles';

const TableComp = (props) => {
	const { importErrors } = props;

	const renderable = importErrors && Array.isArray(importErrors) && importErrors.length

	return (
		<Table>
			<TableHead>
				<TableRow>
					<TableCell align="left">
						Číslo riadku v súbore
					</TableCell>
					<TableCell align="left">
						Popis chyby
					</TableCell>
				</TableRow>
			</TableHead>
			<TableBody>
				{renderable && importErrors.map((e) => (
					<TableRow>
						<TableCell align="left">
							{e.item}
						</TableCell>
						<TableCell align="left">
							{e.reason}
						</TableCell>
					</TableRow>
				))
				}
			</TableBody>
		</Table>
	);
}


TableComp.defaultProps = {
};

export default withStyles(styles)(TableComp);
