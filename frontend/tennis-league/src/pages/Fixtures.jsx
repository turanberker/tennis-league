import React from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';

export default function Fixtures() {
  const { id } = useParams();

  return (
    <Card title="Fikstür">
      <p>Lig ID: {id}</p>
      <p>Burada maç fikstürü listelenecek.</p>
    </Card>
  );
}
