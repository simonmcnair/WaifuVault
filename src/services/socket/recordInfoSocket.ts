import { Nsp, SocketService } from "@tsed/socketio";
import * as SocketIO from "socket.io";
import { RecordInfoPayload } from "../../utils/typeings.js";
import { Inject } from "@tsed/di";
import { ObjectUtils } from "../../utils/Utils.js";
import { FileRepo } from "../../db/repo/FileRepo.js";

@SocketService("/recordInfo")
export class RecordInfoSocket {
    public constructor(@Inject() private repo: FileRepo) {}

    @Nsp
    private nsp: SocketIO.Namespace;

    public async emit(): Promise<boolean> {
        const payload = await this.getPayload();
        return this.nsp.emit("record", payload);
    }

    public async $onConnection(): Promise<void> {
        const payload = await this.getPayload();
        this.nsp.emit("record", payload);
    }

    private async getPayload(): Promise<RecordInfoPayload> {
        const records = await this.repo.getRecordCount();
        const size = (await this.repo.getTotalFileSize()) ?? 0;
        return {
            recordCount: records.toLocaleString(),
            recordSize: ObjectUtils.sizeToHuman(size),
        };
    }
}
