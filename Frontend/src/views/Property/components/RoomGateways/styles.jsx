import { amber, green } from '@material-ui/core/colors';
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
  warning: {
    backgroundColor: amber[700],
    opacity: 0.9,
    marginRight: theme.spacing(1),
    margin: theme.spacing.unit,
  },
  message: {
    display: 'flex',
    alignItems: 'center',
    fontSize: '14px',
  },
  gwItem: {
    margin: theme.spacing.unit,
  },
  demo: {

  },
});
