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
  Button,
  DialogContentText,
} from '@material-ui/core';
import { Delete } from '@material-ui/icons';

import moment from 'moment';
// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';

// Custom components
import { Table as ContentTable, Toolbar } from './components';

import { useStoreState, useStoreActions } from "easy-peasy";

// Component styles
import styles from './style';


function List(props) {
  const { classes } = props;

  const urlBase = '/miestnosti';

  const [selectedRooms, setSelectedRooms] = useState([]);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  //const [search, setSearch] = useState(null);
  const [srchTimeout, setSrchTimeout] = useState(null);
  const [initLoad, setInitLoad] = useState(false);

  const [promptDelete, setProptDelete ] = useState(null);


  const items = useStoreState(state => state.rooms.rooms);
  const loadItems = useStoreActions(actions => actions.rooms.loadAll);
  const isLoading = useStoreState(state => state.rooms.isLoading);
  const error = useStoreState(state => state.rooms.error);
  const firstLoad = useStoreState(state => state.rooms.firstLoadFinished);

  const getItems = (val) => val.rooms;
  const getCount = (val) => val.count;

  /*
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

*/
  const editItem = id => {
    props.history.push(`${urlBase}/${id}`)
  };

  const addNewItem = () => {
    editItem("nova");
  }

  const deleteItem = () => {
    setProptDelete(null);
    //props.deleteRoom(id);
  };

  const performSearch = q => {
    const s = q && q.trim();

    if (s && s.length === 0) {
      setPage(0);
      loadItems();
      //setSearch(null);
    } else {
      setPage(0);
      loadItems(s);
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
      loadItems();
      setInitLoad(true);
    }
  });

  const renderItems = () => {
    const { classes } = props;

    if (isLoading && !firstLoad) {
      return (
        <div className={classes.progressWrapper}>
          <CircularProgress />
        </div>
      );
    }

    if (error) {
      return <Typography variant="h6">{error.codeError}</Typography>;
    }

    if (!items || getCount(items) === 0) {
      //if (!search) {
      return <Typography variant="h4">Zoznam miestností je prázdny</Typography>;
      /*} else {
        return <Typography variant="h4">Neexistuje žiadny záznam pre vyhľadávaný výraz {search}</Typography>;
      }*/
    }

    const keyHeader = {
      name: "Názov",
      value: (itm) => itm.name,
    }

    const valueHeader = [
      {
        name: "Kódové označenie",
        value: i => i.code,
      },
      {
        name: "Posledná aktivita",
        value: i => moment(i.updated_at).format('HH:mm:ss DD.MM.YYYY')
      },
    ]

    const actionHeader = [
      {
        content: "Upraviť",
        callback: i => editItem(i.id),
      },
      {
        content: <Delete />,
        callback: i => setProptDelete(i),
      },
    ]

    return (
      <ContentTable
        //
        onSelect={handleSelect}
        onLimitChange={handleLimitChange}
        onPageChange={handlePageChange}
        items={getItems(items)}
        count={getCount(items)}
        keyHeader={keyHeader}
        valueHeader={valueHeader}
        actionHeader={actionHeader}
        rowsPerPage={limit}
        page={page}
        urlBase={urlBase}
      />
    );
  }

  return (
    <DashboardLayout title="Miestnosti">
      <Dialog
        open={promptDelete}
        onClose={() => { setProptDelete(false) }}
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
          <Button onClick={() => { setProptDelete(false) }} color="primary" autoFocus>
            Zatvoriť
            </Button>
          <Button onClick={() => { deleteItem() }} color="primary">
            Odstrániť
            </Button>
        </DialogActions>
      </Dialog>
      <div className={classes.root}>
        <Toolbar
          selected={selectedRooms}
          add={addNewItem}
          search={performSearch}
          loading={isLoading && initLoad}
          onImport={() => { }}
        />

        <div className={classes.content}>{renderItems()}</div>
      </div>
    </DashboardLayout>
  );
}

export default withStyles(styles)(List);
