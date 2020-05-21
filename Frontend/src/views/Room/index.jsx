import React, { useState, useEffect } from 'react';
import { useHistory } from "react-router-dom";
// Externals
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Grid, CircularProgress } from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';
import { useStoreState, useStoreActions } from 'easy-peasy';

/*
// Custom components
import RoomDetails from './components/RoomDetails';
import RoomGateways from './components/RoomGateways';
*/

import Detail from "./components/Details";


/*
import RoomDetail from '../../components/RoomPicker';
import PropLocation from '../../components/PropLocation';
*/

import Gateways from '../../components/GatewayPicker';
import ActiveUser from '../../components/RoomUser';
// Component styles
const styles = theme => ({
	root: {
		padding: theme.spacing.unit * 4
	}
});

const ItemDetail = (props) => {
	const { classes, match } = props;
	const { params } = match;
	const { id } = params;

	const newIDName = "nova";
	const baseUrl = "/miestnosti"
	const history = useHistory();

	const [ inited, setInited ] = useState(false);
	const [ tabIndex, setTabIndex ] = useState(0);


	const item = useStoreState(state => state.rooms.room);
	const itemPending = useStoreState(state => state.props.propLoading);
	const itemError = useStoreState(state => state.props.propError);

	const loadItem = useStoreActions(actions => actions.rooms.loadRoom);
	const saveItem = useStoreActions(actions => actions.rooms.saveRoom);
	const addItem = useStoreActions(actions => actions.rooms.addRoom);

	const clearItem = useStoreActions(actions => actions.rooms.clearRoom);
	/*
	const loadProp = useStoreActions(actions => actions.props.loadProp);
	
*/
/*
	const saveItem = (room) => {
		if (room) {
			const editedProp = {
				...prop,
				room_id: room.id,
				room: null,
			}
			console.log("save prop", editedProp);
			//saveItem(editedProp);
		}
	}

*/
	const createNewItem = (prop) => {
		if (prop) {
			const editedProp = {
				...prop,
				id: null,
				name: prop.name,
				code: prop.code,
			}
			addItem(editedProp).then(p => {
				console.log("created", p)
				history.push(`${baseUrl}/${p.id}`)
			});
		}
		
	}

	const saveEditedItem = (prop) => {
	if (prop) {
		const editedProp = {
			...prop,
			name: prop.name,
		}
		saveItem(editedProp);
	}
}
	useEffect(() => {
		if (!inited) {
			setInited(true);
			if (id !== newIDName) {
				loadItem(params.id);
			} else {
				clearItem();
			}
		}
	}, []);

	useEffect(() => {
		console.log("item", item);
		console.log("item pen", itemPending);
	}, [item]);

	const windowName = item && item.name ? item.name : `Miestnosť`
	
	const isNew = id == newIDName;

	return (
		<DashboardLayout title={windowName}>
			<div className={classes.root}>
				{itemPending &&
					<div className={classes.progressWrapper}>
						<CircularProgress />
					</div>
				}
				{!itemPending && itemError &&
					<div className={classes.progressWrapper}>
						<h2>Počas načítania nastala chyba</h2>
						{
							itemError && itemError.code === 404 &&
							<h3>
								Miestnosť sa nenašla.
							</h3>
						}
					</div>
				}
				{!itemPending && !itemError &&
					<>
						<Grid
							container
							spacing={4}
						>
							<Grid
								item
								lg={8}
								md={6}
								xl={8}
								xs={12}
							>
								{<Detail
									data={item}
									save={saveEditedItem}
									isNew={isNew}
									createItem={createNewItem}
								/>}
							</Grid>
							<Grid
								item
								lg={4}
								md={6}
								xl={4}
								xs={12}
							>
								{ !isNew && <Gateways 
									id={item.id}
									//save={saveRoom}
								/>}
							</Grid>
						</Grid>
						<Grid
							container
							spacing={4}
						>
							<Grid
								item
								lg={8}
								md={6}
								xl={8}
								xs={12}
							>
								{ !isNew && <ActiveUser 
									roomPending={itemPending}
									room={item}
								/>}
							</Grid>
							<Grid
								item
								lg={4}
								md={6}
								xl={4}
								xs={12}
							>
								{/*<CardDetail  
									id={prop.id}
								/>*/}
							</Grid>
						</Grid>
					</>
				}
			</div>
		</DashboardLayout>
	);
}

ItemDetail.propTypes = {
	classes: PropTypes.object.isRequired
};

export default withStyles(styles)(ItemDetail);
