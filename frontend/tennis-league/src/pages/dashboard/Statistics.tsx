import { Card } from "primereact/card";
import { useEffect, useState } from "react";
import { getStatistics } from "../../api/dashboardService";
import { PlayerStatisticsResponse } from "../../model/dashboard.model";

interface StatisticsProps extends DashboardProps {

}

export default function StatisticsCard({ className = "col-12 md:col-4" }: StatisticsProps) {


    const [statistics, setStatistics] = useState<PlayerStatisticsResponse>();


    useEffect(() => {

        const getStatisticsData = async () => {
            const res = await getStatistics({ limit: 5 })
            if (res) {
                if (res.singlePoints === undefined || res.doublePoints === undefined) {
                    setStatistics(undefined);
                    return;
                }
                setStatistics(res);
            }

        }

        getStatisticsData();

    }, [])




    return (
        <div className={className}>
            <Card title="İstatistikler" style={{ height: '100%' }}>
                {statistics ? (<><div className="flex justify-content-between mb-2">
                    <span>Tekler Puanı</span>
                    <span className="font-bold text-green-500">{statistics?.singlePoints}</span>
                </div>
                    <div className="flex justify-content-between mb-2">
                        <span>Çiftler Puanı</span>
                        <span className="font-bold text-green-500">{statistics?.doublePoints}</span>
                    </div>
                    <div className="flex justify-content-between mb-2">
                        <span>Son 5 maçta kazanılan Tekler Puanı</span>
                        <span className="font-bold text-orange-500">{statistics?.earnedSinglePoints}</span>
                    </div>
                    <div className="flex justify-content-between mb-2">
                        <span>Son 5 maçta kazanılan Çiftler Puanı</span>
                        <span className="font-bold text-orange-500">{statistics?.earnedDoublePoints}</span>
                    </div></>) : (<span>Oyuncu Kaydı Bulunamadı</span>)}


            </Card>
        </div>
    );
}