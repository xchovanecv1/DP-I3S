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
	const { classes, className, importHeader, headerLoading, fileLoading, fileHeader, handleChange, onImportHeader, importActions } = props;

	const rootClassName = classNames(classes.root, className);

	const [importHead, setImportHead] = useState({});
	const [importHeadIncluded, setImportHeadIncluded] = useState(true);

	React.useEffect(() => {
		setImportHead(importHeader);
	}, [importHeader]);

	React.useEffect(() => {
		onImportHeader(importHeadIncluded);
	}, [importHeadIncluded]);

	const renderTable = () => {

		const handleActionSelect = (event, headerKey) => {
			const val = importHead[headerKey];
			val.op_code = (event.target.value);
			setImportHead({
				...importHead,
				[headerKey]: val,
			});

			handleChange(importHeader);
		};
		const handleSelect = (event, headerKey) => {
			const val = importHead[headerKey];
			val.fileKey = (event.target.value);
			setImportHead({
				...importHead,
				[headerKey]: val,
			});

			handleChange(importHeader);
		};

		const actionSelect = (actions, headerKey) => {
			return (
				<FormControl className={classes.formControl}>
					<Select
						className={classes.colSelect}
						labelId="demo-simple-select-placeholder-label-label"
						id="demo-simple-select-placeholder-label"
						value={importHeader[headerKey].op_code}
						onChange={e => handleActionSelect(e, headerKey)}
						displayEmpty
						autoWidth
					//className={classes.selectEmpty}
					>
						{
							actions.map((h, i) => (
								<MenuItem value={h.code}>{h.name}</MenuItem>
							))
						}
					</Select>
					<FormHelperText></FormHelperText>
				</FormControl>
			);
		}

		const select = (headerKey, headerItem) => {
			const len = fileHeader && Array.isArray(fileHeader) && fileHeader.length;
			return (
				<FormControl className={classes.formControl}>
					<Select
						className={classes.colSelect}
						labelId="demo-simple-select-placeholder-label-label"
						id="demo-simple-select-placeholder-label"
						value={importHeader[headerKey].fileKey}
						onChange={e => handleSelect(e, headerKey)}
						displayEmpty
						autoWidth
					//className={classes.selectEmpty}
					>
						{!len &&
							<MenuItem disabled value="">
								<em>Prázdne</em>
							</MenuItem>
						}
						{
							len && fileHeader.map((h, i) => (
								<MenuItem value={i}>{h}</MenuItem>
							)
							)
						}
					</Select>
					<FormHelperText></FormHelperText>
				</FormControl>
			);
		}

		const row = (value, key) => {
			const showActions = importActions && importActions[value.code] && Array.isArray(importActions[value.code]) && importActions[value.code].length;
			console.log(importActions, value.code, importActions[value.code])
			return (
				<TableRow>
					<TableCell align="left">{value.name} ({value.type})</TableCell>
					<TableCell align="right">
						{select(key, value)}
					</TableCell>
					<TableCell align="right">
						{
							showActions &&
							actionSelect(importActions[value.code], key)
						}
					</TableCell>
				</TableRow>);
		}

		const rows = () => {
			const ret = importHeader.map((value, key) => row(value,key));
			ret.unshift(
				<TableRow>
					<TableCell align="left">Názov stĺpca</TableCell>
					<TableCell align="left">Stĺpec súboru</TableCell>
					<TableCell align="left">Akcia</TableCell>
				</TableRow>
			);
			return ret;
		}

		return (
			<Table>
				<TableBody>
					<TableRow>
						<TableCell colspan="3" align="left">
						<FormControlLabel
							control={
							<Checkbox checked={importHeadIncluded} onChange={(e) => { setImportHeadIncluded(e.target.checked) }} value="checkedA" />
							}
							label="Importovaný súbor obsahuje hlavičku na prvom riadku"
						/>
						</TableCell>
					</TableRow>
					{
						rows()
					}
				</TableBody>
			</Table>
		);
	}

	return (
		<Portlet className={rootClassName}>
			<PortletContent noPadding>
				{(headerLoading || fileLoading) && (
					<div className={classes.progressWrapper}>
						<CircularProgress />
					</div>
				)}
				{!(headerLoading || fileLoading) && renderTable(importHead, fileHeader)}
			</PortletContent>
		</Portlet>
	);
}


TableComp.defaultProps = {
	onImportHeader: () => {},
};

export default withStyles(styles)(TableComp);
