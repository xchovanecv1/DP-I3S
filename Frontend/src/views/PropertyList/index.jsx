import React, { useState } from 'react';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import {
  CircularProgress, Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button
} from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';

// Custom components
import { Table as ContentTable, Toolbar } from './components';

import { useStoreState, useStoreActions } from "easy-peasy";

// Component styles
import styles from './style';

import ImportToolbar from "./components/ImportToolbar";
import ImportTable from "./components/Import";
import ErrorTable from "./components/ImportErrors";
import Uploader from "./components/Upload";
import { DropzoneDialog } from 'material-ui-dropzone'

function List(props) {
  const { classes } = props;

  const [selectedRooms, setSelectedRooms] = useState([]);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  //const [search, setSearch] = useState(null);
  const [srchTimeout, setSrchTimeout] = useState(null);
  const [initLoad, setInitLoad] = useState(false);

  const [uploadDialog, setUploadDialog] = useState(false);
  const [importHeaderIncluded, setImportHeaderIncluded] = useState(true);


  const items = useStoreState(state => state.props.props);
  const loadProps = useStoreActions(actions => actions.props.loadAll);
  const isLoading = useStoreState(state => state.props.isLoading);
  const firstLoaded = useStoreState(state => state.props.initalLoaded);

  const importFile = useStoreState(state => state.props.importFile);
  const setImportFile = useStoreActions(actions => actions.props.setImportFile);


  const importHeader = useStoreState(state => state.props.importHeader);
  const headerLoading = useStoreState(state => state.props.headerLoading);
  const loadHeader = useStoreActions(actions => actions.props.loadImportHeader);
  const setImportHeader = useStoreActions(actions => actions.props.setImportHeader);

  const loadFileHeader = useStoreActions(actions => actions.props.loadFileHeader);
  const fileLoading = useStoreState(state => state.props.fileLoading);
  const fileHeader = useStoreState(state => state.props.fileHeader);

  const importActions = useStoreState(state => state.import.actionList);
  const importActionsLoading = useStoreState(state => state.import.actionLoading);
  const loadImportActions = useStoreActions(actions => actions.import.loadActionList);


  const importFileSubmit = useStoreActions(actions => actions.props.importFile);
  const setImportReturn = useStoreActions(actions => actions.props.setImportReturn);
  const importReturn = useStoreState(state => state.props.importReturn);
  const importPending = useStoreState(state => state.props.importLoading);

  const commitPending = useStoreState(state => state.props.commitPending);
  const commitPropImport = useStoreActions(actions => actions.props.commitImport);

  const commitImport = () => {
    commitPropImport(importReturn.valid).then(() => {
      loadProps();
    });
  }

  const submitImport = () => {
    importFileSubmit({
      id: importFile.id,
      head: importHeaderIncluded,
      map: importHeader,
    });
  }

  const handleImportChange = (header) => {
    setImportHeader(header);
  }


  const editRoom = id => {
    props.history.push(`/zariadenia/${id}`)
  };

  const addNewRoom = () => {
    editRoom("nove");
  }

  const deleteRoom = id => {
    //props.deleteRoom(id);
  };

  const performSearch = q => {
    const s = q && q.trim();

    if (s && s.length === 0) {
      setPage(0);
      loadProps();
      //setSearch(null);
    } else {
      setPage(0);
      loadProps(s);
      //setSearch(s);
    }

    clearTimeout(srchTimeout);

    //setSrchTimeout(tmt);
  }

  const handleSelect = selectedRooms => {
    setSelectedRooms(selectedRooms);
  };

  const handleLimitChange = limit => {
    setLimit(limit);
    setPage(0);
  };
  const handlePageChange = page => {
    setPage(page);
  };
/*
  React.useEffect(() => {
    getRooms();
  }, [limit, page]);
*/
  React.useEffect(() => {
    if (!initLoad) {
      loadProps();
      setInitLoad(true);
    }
  });

  React.useEffect(() => {
    if (importFile != null && importFile.id) {
      loadHeader();
      loadImportActions();
      loadFileHeader(importFile.id);
    }
  }, [importFile]);

  const renderImport = () => {
    const { classes } = props;

    if (importPending) {
      return (
        <div className={classes.progressWrapper}>
          <CircularProgress />
        </div>
      );
    }
/*
    if (importErrors) {
      return <Typography variant="h6">{error}</Typography>;
    }

    if (!items || items.count === 0) {
      if (!search) {
        return <Typography variant="h4">Zoznam miestností je prázdny</Typography>;
      } else {
        return <Typography variant="h4">Neexistuje žiadny záznam pre vyhľadávaný výraz {search}</Typography>;
      }
    }
*/

    return (
      (
        <>
          <Typography variant="h4">Importovateľné položky</Typography>
          <ContentTable
            //
            // onSelect={handleSelect}
            //onLimitChange={handleLimitChange}
            // onPageChange={handlePageChange}
            pagination={false}
            readonly={true}
            rooms={importReturn && importReturn.valid}
            count={importReturn && importReturn.valid.length}
          //rowsPerPage={limit}
          //edit={editRoom}
          //delete={deleteRoom}
          //page={page}
          />
          {importReturn && importReturn.invalid_rows && Array.isArray(importReturn.invalid_rows) &&
            <>
              <Typography variant="h4">Položky obsahujúce chyby</Typography>
              <ErrorTable
                importErrors={importReturn.invalid_rows}
              />
            </>
          }
        </>
      )
    );
  }

  const renderRooms = () => {
    const { classes } = props;

    const pending = false;
    const rooms = []
    const error = null

    console.log(firstLoaded);
    if (isLoading && !firstLoaded) {
      return (
        <div className={classes.progressWrapper}>
          <CircularProgress />
        </div>
      );
    }

    if (error) {
      return <Typography variant="h6">{error}</Typography>;
    }

    if (!isLoading && !items || items.count === 0) {
      //if (!search) {
        return <Typography variant="h4">Zoznam zariadení je prázdny</Typography>;
      /*} else {
        return <Typography variant="h4">Neexistuje žiadny záznam pre vyhľadávaný výraz {search}</Typography>;
      }*/
    }

    return (
      <ContentTable
        //
        onSelect={handleSelect}
        onLimitChange={handleLimitChange}
        onPageChange={handlePageChange}
        rooms={items.properties}
        count={items.count}
        rowsPerPage={limit}
        edit={editRoom}
        delete={deleteRoom}
        page={page}
      />
    );
  }

  return (
    <DashboardLayout title="Zariadenia">
      <div className={classes.root}>
        {importReturn != null &&
        <>
          <ImportToolbar
            close={() => { setImportReturn(null) }}
            confirm={() => { commitImport() }}
          />
          <div className={classes.content}>{renderImport()}</div>
        </>
        }
        {importReturn == null &&
        <>
          <Toolbar
            selected={selectedRooms}
            add={addNewRoom}
            search={performSearch}
            loading={isLoading && initLoad}
            onImport={() => { setUploadDialog(true) }}
          />

          <div className={classes.content}>{renderRooms()}</div>
          
        </>
        }
        <Dialog
          open={uploadDialog}
          onClose={() => setUploadDialog(false)}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogContent>
            {/* (headerLoading || fileLoading) &&
              <div className={classes.progressWrapper}>
                <CircularProgress />
              </div>*/
            }
            <Uploader
              fileUploaded={(f) => {
                setImportFile(f);
                setUploadDialog(false);
              }
              }
            />
            {/*<ImportTableImportTable 
                importHeader={importHeader}
                headerLoading={headerLoading}
                fileLoading={fileLoading}
                fileHeader={fileHeader}
                handleChange={handleImportChange}
              />*/}
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => { setImportFile(null) }}
              color="primary" autoFocus
            >
              Zatvoriť
            </Button>
          </DialogActions>
        </Dialog>
        <Dialog
          open={importFile != null && importReturn == null}
          //onClose={dialog.close}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogContent>
            {(headerLoading || fileLoading || importActionsLoading) &&
              <div className={classes.progressWrapper}>
                <CircularProgress />
              </div>
            }
            {!(headerLoading || fileLoading || importActionsLoading) &&
              <ImportTable
                importHeader={importHeader}
                headerLoading={headerLoading}
                fileLoading={fileLoading}
                fileHeader={fileHeader}
                handleChange={handleImportChange}
                onImportHeader={setImportHeaderIncluded}
                importActions={importActions}
              />
            }
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => { setImportFile(null) }}
              color="secondary" autoFocus
            >
              Zatvoriť
            </Button>
            <Button
              onClick={() => { submitImport() }}
              color="primary" autoFocus
            >
              Importovať
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    </DashboardLayout>
  );
}

export default withStyles(styles)(List);
