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
	const { classes, className, rooms, count, rowsPerPage, page, pagination, readonly } = props;

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
		const { rooms, onSelect } = props;

		let selected;

		if (event.target.checked) {
			selected = rooms.map(data => data.id);
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
				<Dialog
					open={dialog.open}
					onClose={dialog.close}
					aria-labelledby="alert-dialog-title"
					aria-describedby="alert-dialog-description"
				>
					<DialogTitle id="alert-dialog-title">{"Chcete odstrániť zvolenú miestnosť?"}</DialogTitle>
					<DialogContent>
						<DialogContentText id="alert-dialog-description">
							Operácia slúži pre odstránenie požadovanej miestnosti. Tento proces je nezvratný.
            </DialogContentText>
					</DialogContent>
					<DialogActions>
						<Button onClick={dialog.close} color="primary" autoFocus>
							Zatvoriť
            </Button>
						<Button onClick={dialog.agree} color="primary">
							Odstrániť
            </Button>
					</DialogActions>
				</Dialog>
				<PerfectScrollbar>
					<Table>
						<TableHead>
							<TableRow>
								<TableCell align="left">
									{!readonly &&
										<Checkbox
											checked={selected.length === rooms.length}
											color="primary"
											indeterminate={
												selected.length > 0 &&
												selected.length < rooms.length
											}
											onChange={handleSelectAll}
										/>
									}
									Názov
                </TableCell>
								<TableCell align="left">Kódové označnie</TableCell>
								<TableCell align="left">Cena</TableCell>
								<TableCell align="left">Dátum zadováženia</TableCell>
								<TableCell align="left">Miestnosť</TableCell>
								{!readonly && <TableCell align="left"></TableCell> }
							</TableRow>
						</TableHead>
						<TableBody>
							{Array.isArray(rooms) && rooms
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
													{getInitials(data.name)}
												</Avatar>
												{!readonly && <Link to={`/zariadenia/${data.id}`}>
													<Typography
														className={classes.nameText}
														variant="body1"
													>
														{data.name}
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
										<TableCell className={classes.tableCell}>
											{data.code}
										</TableCell>
										<TableCell className={classes.tableCell}>
											{data.price}€
										</TableCell>
										<TableCell className={classes.tableCell}>
											{moment(data.acquired_at).format('DD.MM.YYYY')}
										</TableCell>
										<TableCell className={classes.tableCell}>
											{data.room && data.room.name}
										</TableCell>
										{!readonly &&
										<TableCell className={classes.tableCell}>
											<Button
												color="primary"
												size="small"
												variant="outlined"
												onClick={() => { props.edit(data.id) }}
												className={classes.button}
											>
												Upraviť
                      						</Button>
											<Button
												color="primary"
												size="small"
												variant="outlined"
												onClick={() => { deleteRoom(data.id) }}
												className={classes.button}
											>
												<Delete />
											</Button>
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
	rooms: [],
	count: 0,
	onSelect: () => { },
	edit: (e) => { console.log(e) },
	onShowDetails: () => { },
	pagination: true,
	readonly: false,
};

export default withStyles(styles)(TableComp);
