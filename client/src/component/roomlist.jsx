import React, { useState, useEffect } from 'react';
import axios from 'axios';

const RoomList = () => {
  const [rooms, setRooms] = useState([]);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const response = await axios.get('http://localhost:8001/room');
        setRooms(response.data);
      } catch (err) {
        console.error(err);
      }
    };

    fetchRooms();
  }, []);

  return (
    <div className="room-list">
      <h2>Available Rooms</h2>
      <ul>
        {rooms.map((room) => (
          <li key={room.id}>
            <h3>{room.name}</h3>
            <p>{room.description}</p>
            <button>Reserve</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default RoomList;
