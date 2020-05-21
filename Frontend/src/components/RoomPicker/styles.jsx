export default theme => ({
  root: {},
  field: {
    margin: theme.spacing.unit * 3
  },
  textField: {
    width: '420px',
    maxWidth: '100%',
    marginRight: theme.spacing.unit * 3
  },
  portletFooter: {
    paddingLeft: theme.spacing.unit * 3,
    paddingRight: theme.spacing.unit * 3,
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2
  },
  button: {
    margin: theme.spacing.unit,
  },
  message: {
    display: 'flex',
    alignItems: 'center',
    fontSize: '16px',
  },
  warningIco: {
    color: theme.palette.warning.main,
    marginRight: theme.spacing.unit,
  }
});
