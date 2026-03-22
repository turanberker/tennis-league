import { Card } from "primereact/card";
import { useLeague } from "../hooks/useLeague";
import { DataTable } from "primereact/datatable";
import { Button } from "primereact/button";
import { Player, Sex } from "../model/player.model";
import { CreateTeamRequest, TeamResponse } from "../model/team.model";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import * as yup from 'yup';
import { yupResolver } from "@hookform/resolvers/yup";
import { getPlayers } from "../api/playersService";
import { createTeam, getTeams } from "../api/leagueService";
import { Dialog } from "primereact/dialog";
import { Column } from "primereact/column";
import { Dropdown } from "primereact/dropdown";
import { InputText } from "primereact/inputtext";
import { Sidebar } from "primereact/sidebar";
import FormItem from "./FormItem";
import { isFieldRequired } from "../helper/form.helper";
import { LEAGUE_CATEGORY, LEAGUE_FORMAT } from "../model/league.model";
import { useAsyncError } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";

interface LeagueTeamsProps {
    leagueId: string;
}
type CreateTeamForm = {
    name: string;
    player1: Player | null;
    player2: Player | null;
};

const schema = yup.object({
    name: yup
        .string()
        .required('Takım adı zorunludur')
        .min(5, 'Takım adı en az 5 karakter olmalı')
        .max(75, 'Takım adı en fazla 75 karakter olabilir'),
    player1: yup.mixed<Player>().nullable().required('Birinci oyuncuyu seçin'),
    player2: yup
        .mixed<Player>()
        .nullable()
        .required('İkinci oyuncuyu seçin')
        .test(
            'different-player',
            'İki oyuncu birbirinden farklı olmalıdır',
            function (value) {
                const { player1 } = this.parent;
                if (!value || !player1) return true;
                return (value as Player).id !== (player1 as Player).id;
            },
        ),
}) as yup.Schema<CreateTeamForm>; // TypeScript tip uyumu için

export const LeagueTeams: React.FC<LeagueTeamsProps> = ({ leagueId }) => {

    const { data: league, updateLeagueCache } = useLeague(leagueId)
    const [teams, setTeams] = useState<TeamResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);

    const [createDialogVisible, setCreateDialogVisible] =
        useState<boolean>(false);

    const methods
        = useForm<CreateTeamForm>({
            resolver: yupResolver(schema as any),
            defaultValues: { name: '', player1: null, player2: null },
        });
    const [playersLoaded, setPlayerLoaded] = useState<Boolean>(false);
    const [playerList1, setPlayerList1] = useState<Player[]>();
    const [playerList2, setPlayerList2] = useState<Player[]>();


    const loadPlayers = async () => {
        if (!playersLoaded) {
            if (league?.category === LEAGUE_CATEGORY.FEMALE) {
                const res = await getPlayers({ sex: Sex.Female });
                setPlayerList1(res)
                setPlayerList2(res)
            }
            else if (league?.category === LEAGUE_CATEGORY.MALE) {
                const res = await getPlayers({ sex: Sex.Male });
                setPlayerList1(res)
                setPlayerList2(res)
            } else if (league?.category === LEAGUE_CATEGORY.MIX) {


                const res1 = await getPlayers({ sex: Sex.Male });
                setPlayerList1(res1)

                const res2 = await getPlayers({ sex: Sex.Female });
                setPlayerList2(res2)

            }
        }
        setPlayerLoaded(true)

    }

    // Takım oluştur
    const onSubmit = async (data: CreateTeamForm) => {
        if (!leagueId) return;

        const payload: CreateTeamRequest = {
            name: data.name,
            playerIds: [data.player1!.id, data.player2!.id],
        };
        const res: { teamId: String, totalAttendanceCount: number } = await createTeam(leagueId, payload);
        if (res) {
            updateLeagueCache({
                totalAttentance: res.totalAttendanceCount
            });

            setCreateDialogVisible(false);
            methods.reset();
            loadTeams();
        }

    };

    // Lig takımlarını yükle
    const loadTeams = async (): Promise<void> => {
        if (!leagueId) return;

        setLoading(true);
        const res: TeamResponse[] = await getTeams(leagueId);
        setTeams(res);
        setLoading(false);
    };

    useEffect(() => {
        loadTeams();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [leagueId]);

    const playerLabelItemTemplate = (option: Player) => {
        return (
            option ? option.name + ' ' + option.surname : 'Oyuncu seçin'
        ) as string;
    };

    const header = () => {
        return (
            <div className="flex justify-content-end">
                <Button
                    label="Yeni Takım"
                    icon="pi pi-plus"
                    size="small"
                    onClick={() => {
                        loadPlayers();
                        setCreateDialogVisible(true)
                    }}
                />
            </div>
        );
    };

    return (
        <>
            <Card
                title="Takımlar"
            >

                <DataTable
                    value={teams}
                    loading={loading}
                    emptyMessage="Takım bulunamadı"
                    tableStyle={{ minWidth: '50rem' }}
                    header={header}
                    key="id"
                >
                    <Column field="name" header="Takım Adı" />
                </DataTable>
            </Card>

            {/* Yeni Takım Dialog */}
            <Sidebar
                header="Yeni Takım Oluştur"
                visible={createDialogVisible}
                className="w-full md:w-25rem"
                position="right"
                onHide={() => setCreateDialogVisible(false)}
            >

                <FormProvider {...methods}>
                    <form onSubmit={methods.handleSubmit(onSubmit)} className="p-fluid">
                        <FormItem label="Takım Adı"
                            name="name"
                            required={isFieldRequired(schema, "name")}>
                            <InputText
                                {...methods.register('name')}
                                className={methods.formState.errors.name ? 'p-invalid' : ''}
                                placeholder="Takım adı girin"
                            />

                        </FormItem>
                        <FormItem label={league?.category === LEAGUE_CATEGORY.MIX ? "Erkek Oyuncu" : "1.Oyuncu"}
                            name="player1"
                            required={isFieldRequired(schema, "player1")}>
                            <Controller
                                name="player1"
                                control={methods.control}
                                render={({ field }) => (
                                    <Dropdown
                                        {...field}
                                        onChange={(e) => field.onChange(e.value)}
                                        filterMatchMode="contains"
                                        filter
                                        filterBy="name,surname"
                                        filterLocale="tr"
                                        options={playerList1}
                                        dataKey="id"
                                        itemTemplate={playerLabelItemTemplate}
                                        valueTemplate={playerLabelItemTemplate}
                                        placeholder="Oyuncu 1 seçin"
                                        className={methods.formState.errors.player1 ? 'p-invalid' : ''}
                                    />
                                )}
                            />
                        </FormItem>

                        <FormItem label={league?.category === LEAGUE_CATEGORY.MIX ? "Kadın Oyuncu" : "2.Oyuncu"}
                            name="player1"
                            required={isFieldRequired(schema, "player2")}>
                            <Controller
                                name="player2"
                                control={methods.control}
                                render={({ field }) => (
                                    <Dropdown
                                        {...field}
                                        onChange={(e) => field.onChange(e.value)}
                                        filterMatchMode="contains"
                                        filter
                                        filterBy="name,surname"
                                        filterLocale="tr"
                                        options={playerList2}
                                        dataKey="id"
                                        itemTemplate={playerLabelItemTemplate}
                                        valueTemplate={playerLabelItemTemplate}
                                        placeholder="Oyuncu 1 seçin"
                                        className={methods.formState.errors.player2 ? 'p-invalid' : ''}
                                    />
                                )}
                            />
                        </FormItem>
                        <div className="mt-4 flex gap-2">
                            <Button
                                type="button"
                                label="İptal"
                                icon="pi pi-times"
                                outlined
                                severity="secondary"
                                onClick={() => setCreateDialogVisible(false)}
                                className="w-full"
                            />
                            <Button
                                type="submit"
                                label="Kaydet"
                                icon="pi pi-check"
                                loading={methods.formState.isSubmitting}
                                className="w-full"
                            />
                        </div>
                    </form>
                </FormProvider>
            </Sidebar>
        </>
    );

}