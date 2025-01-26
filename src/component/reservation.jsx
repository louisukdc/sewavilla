import React, { useState } from 'react';
import axios from 'axios';

const Reservation = ({ roomId }) => {
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [error, setError] = useState('');

  const handleReservation = async (e) => {
    e.preventDefault();

    try {
      const token = localStorage.getItem('token');
      const response = await axios.post(
        'http://localhost:8001/reservation',
        { roomId, startDate, endDate },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      alert('Reservation successful');
    } catch (err) {
      setError('Failed to make reservation');
    }
  };

  return (
    <div className="reservation-container">
      <h2>Reserve Room</h2>
      <form onSubmit={handleReservation}>
        <input
          type="date"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
        />
        <input
          type="date"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
        />
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Make Reservation</button>
      </form>
    </div>
  );
};

export default Reservation;
