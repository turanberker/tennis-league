import React from "react";
import { useParams } from "react-router-dom";
import { Card } from "primereact/card";

export default function Standings() {
  const { id } = useParams();

  return (
    <Card title="Puan Durumu">
      <p>Lig ID: {id}</p>
      <p>Burada puan durumu listelenecek.</p>
    </Card>
  );
}