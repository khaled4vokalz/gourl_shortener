import React, { useState } from "react";
import { createShortUrl } from "../api/service";

const Shortener = () => {
  const [longUrl, setLongUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const result = await createShortUrl(longUrl);
      setShortUrl(result.shortened_url);
    } catch (error) {
      console.error("Error creating short URL:", error);
    }
  };

  return (
    <div className="shortener-container">
      <h1>URL Shortener</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="url"
          placeholder="Enter your long URL"
          value={longUrl}
          onChange={(e) => setLongUrl(e.target.value)}
          required
        />
        <button type="submit">Shorten</button>
      </form>
      {shortUrl && (
        <div className="short-url-container">
          <p>Short URL:</p>
          <a href={shortUrl} target="_blank" rel="noopener noreferrer">
            {shortUrl}
          </a>
        </div>
      )}
    </div>
  );
};

export default Shortener;
