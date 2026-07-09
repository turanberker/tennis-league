import {Card} from "primereact/card";
import {addPlayer, getPlayers} from "../api/mail/leagueService";
import {useCallback, useEffect, useState} from "react";
import {useLeague} from "../hooks/useLeague";
import {LEAGUE_CATEGORY, SingleLeagueAttendanceResponse} from "../model/league.model";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import Guard from "../helper/Guard";
import {Role} from "../model/user.model";
import {Button} from "primereact/button";
import {Player, Sex} from "../model/player.model";
import * as yup from "yup";
import {Controller, FormProvider, useForm} from "react-hook-form";
import {yupResolver} from "@hookform/resolvers/yup";
import {getPlayers as getAvailablePlayers} from "../api/user/playersService";
import {Sidebar} from "primereact/sidebar";
import FormItem from "./FormItem";
import {isFieldRequired} from "../helper/form.helper";
import {Dropdown} from "primereact/dropdown";

interface LeaguePlayersProps {
    leagueId: string;
}


const schema = yup.object({
    player: yup.mixed<string>().required('Oyuncuyu seçin')
});

type CreatePlayerForm = yup.InferType<typeof schema>;

export const LeaguePlayers: React.FC<LeaguePlayersProps> = ({leagueId}) => {
    const {data: league, updateLeagueCache} = useLeague(leagueId)
    const [players, setPlayers] = useState<SingleLeagueAttendanceResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [createDialogVisible, setCreateDialogVisible] = useState<boolean>(false);
    const [availablePlayers, setAvailablePlayers] = useState<Player[]>();
    const [playersLoaded, setPlayerLoaded] = useState<Boolean>(false);

    const methods
        = useForm<CreatePlayerForm>({
        resolver: yupResolver(schema as any),
        defaultValues: {player: undefined},
    });

    const onSubmit = async (data: CreatePlayerForm) => {
        if (!leagueId) return;


        const res: { playerId: string, totalAttendanceCount: number } = await addPlayer(leagueId, data.player);
        if (res) {
            updateLeagueCache({
                totalAttentance: res.totalAttendanceCount
            });

            setCreateDialogVisible(false);
            methods.reset();
            loadPlayers();
        }

    };

    const loadAvailablePlayers = async () => {
        if (!playersLoaded) {
            const res = await getAvailablePlayers({sex: league?.category === LEAGUE_CATEGORY.FEMALE ? Sex.Female : Sex.Male});
            setAvailablePlayers(res)
        }
        setPlayerLoaded(true)

    }

    const loadPlayers = useCallback(async () => {
        if (!leagueId) return;

        setLoading(true);
        const res = await getPlayers(leagueId);
        setPlayers(res);
        setLoading(false);
    }, [leagueId]);

    useEffect(() => {
        loadPlayers();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [leagueId]);

    const header = () => {
        return (
            <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
                <div className="flex justify-content-end">
                    <Button
                        label="Yeni Oyuncu"
                        icon="pi pi-plus"
                        size="small"
                        onClick={() => {
                            loadAvailablePlayers();
                            setCreateDialogVisible(true)
                        }}
                    />
                </div>
            </Guard>
        );
    };
    const playerLabelItemTemplate = (option: Player) => {
        return (
            option ? option.name + ' ' + option.surname : 'Oyuncu seçin'
        ) as string;
    };
    return (<>
            <Card
                title="Oyuncular"
            >
                <DataTable
                    value={players}
                    loading={loading}
                    emptyMessage="Oyuncu bulunamadı"
                    tableStyle={{minWidth: '50rem'}}
                    header={header}
                    key="id"
                >
                    <Column field="firstname" header="Adı"/>
                    <Column field="surname" header="Soyadı"/>
                    <Column field="power" header="Gücü"/>
                </DataTable>

            </Card>


            {/* Yeni Takım Dialog */}
            <Sidebar
                header="Yeni Oyuncu Ekle"
                visible={createDialogVisible}
                className="w-full md:w-25rem"
                position="right"
                onHide={() => setCreateDialogVisible(false)}
            >

                <FormProvider {...methods}>
                    <form onSubmit={methods.handleSubmit(onSubmit)} className="p-fluid">

                        <FormItem label={"Oyuncu"}
                                  name="player"
                                  required={isFieldRequired(schema, "player")}>
                            <Controller
                                name="player"
                                control={methods.control}
                                render={({field}) => (
                                    <Dropdown
                                        {...field}
                                        value={field.value}
                                        onChange={(e) => field.onChange(e.value)}
                                        filterMatchMode="contains"
                                        filter
                                        filterBy="name,surname"
                                        filterLocale="tr"
                                        optionLabel="name"
                                        options={availablePlayers}
                                        optionValue="id"
                                        dataKey="id"
                                        itemTemplate={playerLabelItemTemplate}
                                        valueTemplate={playerLabelItemTemplate}
                                        placeholder="Oyuncu 1 seçin"
                                        className={methods.formState.errors.player ? 'p-invalid' : ''}
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