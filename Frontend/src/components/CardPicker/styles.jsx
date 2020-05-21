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
  },
  portletContent: {
    paddingTop: '0'
  },
  portletFooter: {
    display: 'flex',
    justifyContent: 'flex-end'
  },
  product: {
    padding: theme.spacing.unit,
    display: 'flex',
    alignItems: 'center',
    '&:not(:first-of-type)': {
      borderTop: `1px solid ${theme.palette.divider}`
    }
  },
  productImageWrapper: {
    borderRadius: '5px',
    overflow: 'hidden',
    height: '64px',
    width: '64px'
  },
  productImage: {
    width: '100%',
    height: 'auto'
  },
  productDetails: {
    marginLeft: theme.spacing.unit * 2,
    flexGrow: 1
  },
  productTitle: {},
  productTimestamp: {
    marginTop: theme.spacing.unit,
    color: theme.palette.text.secondary
  }
});
