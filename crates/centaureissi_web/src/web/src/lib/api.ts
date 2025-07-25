import { token } from "./stores.js";
import { browser } from "$app/environment";
import axios from "axios";

let url = (browser && window.location.origin) || "";
let current_token: string;
token.subscribe((value) => {
  current_token = value;
});

export async function login({ username, password }) {
  const endpoint = url + "/api/auth/login";
  const response = await axios({
    url: endpoint,
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: JSON.stringify({
      username,
      password,
    }),
  });
  token.set(response.data.token);
  return response.data;
}

export async function register({ email, username, password }) {
  const endpoint = url + "/api/user/register";
  const response = await axios({
    url: endpoint,
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: JSON.stringify({
      email,
      username,
      password,
    }),
  });
  token.set(response.data.token);
  return response.data;
}

export async function updateToken() {
  const endpoint = url + "/api/auth/token";
  const response = await axios({
    url: endpoint,
    method: "POST",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
  });
  token.set(response.data.token);
  return response.data;
}

export async function getUserProfile() {
  const endpoint = url + "/api/user/profile";
  const response = await axios({
    url: endpoint,
    method: "GET",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    withCredentials: true,
  });
  console.log(response.data);
  return response.data;
}

export async function updateUserProfile({ email, username }) {
  const endpoint = url + "/api/user/update";
  const response = await axios({
    url: endpoint,
    method: "PUT",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    withCredentials: true,
    data: {
      email,
      username,
    },
  });
  return response.data;
}
export async function updateUserPassword({ old_password, new_password }) {
  const endpoint = url + "/api/user/update-password";
  const response = await axios({
    url: endpoint,
    method: "PUT",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    withCredentials: true,
    data: {
      old_password,
      new_password,
    },
  });
  return response.data;
}

export async function search({
  q,
  page,
  perPage,
}: {
  q: any;
  page: any;
  perPage?: any;
}) {
  if (!q) {
    q = "";
  }
  if (!perPage) {
    perPage = 10;
  }
  let endpoint =
    url + "/api/search?q=" + q + "&page=" + page + "&per_page=" + perPage;

  const response = await axios({
    url: endpoint,
    method: "GET",
    headers: {
      Authorization: "Bearer " + current_token,
    },
    withCredentials: true,
  });
  return response.data;
}

export async function stats() {
  let endpoint = url + "/api/stats";

  const response = await axios({
    url: endpoint,
    method: "GET",
    headers: {
      Authorization: "Bearer " + current_token,
    },
    withCredentials: true,
  });
  return response.data;
}

export async function getMessageDetail(message_id: string): any {
  const endpoint = url + `/api/messages/${message_id}/info`;
  const response = await axios({
    url: endpoint,
    method: "GET",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    withCredentials: true,
  });
  return response.data;
}

export async function getMessageByType(
  message_id: string,
  message_type: string,
): any {
  const endpoint = url + `/api/messages/${message_id}/${message_type}`;
  const response = await axios({
    url: endpoint,
    method: "GET",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    withCredentials: true,
  });
  return response.data;
}

export async function tokenSign(path: string) {
  const endpoint = url + "/api/auth/token/sign";
  const response = await axios({
    url: endpoint,
    method: "POST",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    data: JSON.stringify({
      path,
    }),
  });
  return response.data.token;
}

export async function getToken() {
  const endpoint = url + "/api/auth/token";
  const response = await axios({
    url: endpoint,
    method: "POST",
    headers: {
      Authorization: "Bearer " + current_token,
      "Content-Type": "application/json",
    },
    data: JSON.stringify({}),
  });
  return response.data.token;
}
export async function downloadAttachment(
  message_id: string,
  id: number,
  token: string,
) {
  const endpoint =
    url + `/api/messages/${message_id}/attachment/${id}?token=${token}`;
  window.open(endpoint);
}
