import {
  ALL_ROOMS_LOADED,
  ALL_ROOMS_PENDING,
  ALL_ROOMS_FAILED,

  ROOM_LOADED,
  ROOM_PENDING,
  ROOM_FAILED,

  ROOM_SAVED,
  ROOM_SAVE_FAILED,
  ROOM_SAVE_PENDING,

  ROOM_DELETED,
  ROOM_DELETE_FAILED,
  ROOM_DELETE_PENDING,

} from "../actions/rooms"

const initialState = {
  roomsPending: true,
  rooms: null,
  roomsError: null,

  roomPending: true,
  room: null,
  roomError: null,

  saved: false,
  savePending: false,
  saveError: null,

  deleted: false,
  deletePending: false,
  deleteError: null,

}

export default (state = initialState, action) => {
  switch (action.type) {
    case ALL_ROOMS_LOADED:
    case ALL_ROOMS_PENDING:
    case ALL_ROOMS_FAILED:

    case ROOM_LOADED:
    case ROOM_PENDING:
    case ROOM_FAILED:

    case ROOM_SAVED:
    case ROOM_SAVE_PENDING:
    case ROOM_SAVE_FAILED:


    case ROOM_DELETED:
    case ROOM_DELETE_FAILED:
    case ROOM_DELETE_PENDING:

      return {
        ...state,
        ...action,
      }
    default:
      return state
  }
}