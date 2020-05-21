import axios from "../../common/axios"

export const getProfileInfo = () => {
	return axios.get("/api/api/profile");
}
