import React, { useState } from 'react';
import { Link } from 'react-router-dom';

import { useStoreActions } from "easy-peasy"
// Externals
import classNames from 'classnames';
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


import { DropzoneArea } from 'material-ui-dropzone'

// Component styles
import styles from './styles';

const TableComp = (props) => {
	const { classes, className, fileUploaded } = props;

	const rootClassName = classNames(classes.root, className);

	const [files, setFiles] = useState([]);
	const [uploading, setUploading] = useState(false);

	const upload = useStoreActions(actions => actions.import.upload);

	React.useEffect(() => {
		if (files.length > 0) {
			setUploading(true);
			upload(files[0])
				.then(d => {
					fileUploaded(d);
				})
				.catch(() => {
					setUploading(false);
				});
		}
	}, [files])

	return (
		<Portlet className={rootClassName}>
			<PortletContent noPadding>
				{(uploading) &&
					<div className={classes.progressWrapper}>
						<CircularProgress />
					</div>
				}
				{!uploading &&
					<DropzoneArea
						acceptedFiles={[
							'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
						]}
						onChange={setFiles}
						filesLimit={1}
						dropzoneText={"Priložte súbor, alebo ho kliknutím vyberte."}
						showPreviews={false}
						showPreviewsInDropzone={false}
						useChipsForPreview={false}
						showFileNamesInPreview={false}
					/>
				}
			</PortletContent>
		</Portlet>
	);
}


TableComp.defaultProps = {
	fileUploaded: () => {},
};

export default withStyles(styles)(TableComp);
