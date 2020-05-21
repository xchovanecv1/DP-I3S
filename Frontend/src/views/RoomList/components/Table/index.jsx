import React, { useState } from 'react';
import { Link } from 'react-router-dom';

// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';

// Material helpers
import { withStyles } from '@material-ui/core';

// Dialog helpers
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@material-ui/core';

// Material components
import {
	Avatar,
	Checkbox,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
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
import { Delete } from '@material-ui/icons';

const TableComp = (props) => {
	const { classes, className, items, count, rowsPerPage, page, pagination, readonly, keyHeader, valueHeader, urlBase, actionHeader } = props;

	const [activeTab, setActiveTab] = useState(null);
	const [selected, setSelected] = useState([]);

	const [dialog, setDialog] = useState({
		open: false,
		agree: () => { },
		close: () => setDialog({ ...dialog, open: false }),
	});

	const closeDialog = () => {
		setDialog({ ...dialog, open: false });
	}

	const rootClassName = classNames(classes.root, className);

	const handleSelectAll = event => {
		const { items, onSelect } = props;

		let selected;

		if (event.target.checked) {
			selected = items.map(data => data.id);
		} else {
			selected = [];
		}

		setSelected(selected);
		onSelect(selected);
	};

	const deleteRoom = (id) => {
		setDialog({
			...dialog,
			open: true,
			agree: () => {
				props.delete(id)
				closeDialog()
			}
		});
	}

	const handleSelectOne = (event, id) => {
		const { onSelect } = props;

		const selectedIndex = selected.indexOf(id);
		let newselected = [];

		if (selectedIndex === -1) {
			newselected = newselected.concat(selected, id);
		} else if (selectedIndex === 0) {
			newselected = newselected.concat(selected.slice(1));
		} else if (selectedIndex === selected.length - 1) {
			newselected = newselected.concat(selected.slice(0, -1));
		} else if (selectedIndex > 0) {
			newselected = newselected.concat(
				selected.slice(0, selectedIndex),
				selected.slice(selectedIndex + 1)
			);
		}

		setSelected(newselected);
		onSelect(newselected);
	};

	const handleChangePage = (event, page) => props.onPageChange(page);

	const handleChangeRowsPerPage = event => props.onLimitChange(event.target.value);


	return (
		<Portlet className={rootClassName}>
			<PortletContent noPadding>
				<PerfectScrollbar>
					<Table>
						<TableHead>
							<TableRow>
								<TableCell align="left">
									{!readonly &&
										<Checkbox
											checked={selected.length === items.length}
											color="primary"
											indeterminate={
												selected.length > 0 &&
												selected.length < items.length
											}
											onChange={handleSelectAll}
										/>
									}
									{keyHeader.name}
                				</TableCell>
								{
									Array.isArray(valueHeader) && valueHeader.length && valueHeader.map((v) => (
										<TableCell align="left">{ v.name }</TableCell>
									))
								}
								{!readonly && Array.isArray(actionHeader) && actionHeader.length && <TableCell align="right"></TableCell> }
							</TableRow>
						</TableHead>
						<TableBody>
							{Array.isArray(items) && items
								.map(data => (
									<TableRow
										className={classes.tableRow}
										hover
										key={data.id}
										selected={selected.indexOf(data.id) !== -1}
									>
										<TableCell className={classes.tableCell}>
											<div className={classes.tableCellInner}>
												{!readonly &&
													<Checkbox
														checked={selected.indexOf(data.id) !== -1}
														color="primary"
														onChange={event =>
															handleSelectOne(event, data.id)
														}
														value="true"
													/>
												}
												<Avatar
													className={classes.avatar}
													src={data.avatarUrl}
												>
													{getInitials(keyHeader.value(data))}
												</Avatar>
												{!readonly && <Link to={`${urlBase}/${data.id}`}>
													<Typography
														className={classes.nameText}
														variant="body1"
													>
														{keyHeader.value(data)}
													</Typography>
												</Link>
												}
												{readonly &&
													<Typography
														className={classes.nameText}
														variant="body1"
													>
														{data.name}
													</Typography>
												}
											</div>
										</TableCell>
										{
											Array.isArray(valueHeader) && valueHeader.length && valueHeader.map((v) => (
												<TableCell className={classes.tableCell}>
													{ v.value(data) }
												</TableCell>
											))
										}
										{!readonly && Array.isArray(actionHeader) && actionHeader.length && 
										<TableCell className={classes.tableCell}>
											{
												actionHeader.map(a => (
													<Button
														color="primary"
														size="small"
														variant="outlined"
														onClick={() => { a.callback(data) }}
														className={classes.button}
													>
														{a.content}
													</Button>
												))
											}
										</TableCell>
										}
									</TableRow>
								))}
						</TableBody>
					</Table>
				</PerfectScrollbar>
				{pagination &&
					<TablePagination
						backIconButtonProps={{
							'aria-label': 'Previous Page'
						}}
						component="div"
						count={count}
						nextIconButtonProps={{
							'aria-label': 'Next Page'
						}}
						onChangePage={handleChangePage}
						onChangeRowsPerPage={handleChangeRowsPerPage}
						page={page}
						count={count}
						rowsPerPage={rowsPerPage}
						rowsPerPageOptions={[1, 5, 10, 25]}
					/>
				}
			</PortletContent>
		</Portlet>
	);
}


TableComp.defaultProps = {
	items: [],
	count: 0,
	onSelect: () => { },
	edit: (e) => { console.log(e) },
	onShowDetails: () => { },
	pagination: true,
	readonly: false,
	keyHeader: {
		name: "Názov",
		value: (itm) => itm.name,
	},
};

export default withStyles(styles)(TableComp);
