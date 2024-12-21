import axios from "axios";

const API_BASE_URL =
  process.env.REACT_APP_BACKEND_URL || "http://localhost:8082";

export const createShortUrl = async (longUrl) => {
  const response = await axios.post(`${API_BASE_URL}/shorten`, {
    url: longUrl,
  });
  return response.data;
};

export const fetchOriginalUrl = async (shortCode) => {
  const response = await axios.get(`${API_BASE_URL}/${shortCode}`);
  return response.data;
};
