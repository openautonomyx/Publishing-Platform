/**
 * Team Service - Manage team members and permissions
 */

import { AxiosInstance } from 'axios';
import { Team, TeamMember, InviteRequest } from '../types';

export class TeamService {
  constructor(private apiClient: AxiosInstance) {}

  async getTeam(): Promise<Team> {
    const response = await this.apiClient.get<Team>('/team');
    return response.data;
  }

  async invite(request: InviteRequest): Promise<TeamMember> {
    const response = await this.apiClient.post<TeamMember>(
      '/team/invite',
      request
    );
    return response.data;
  }

  async removeMember(memberId: string): Promise<void> {
    await this.apiClient.delete(`/team/members/${memberId}`);
  }

  async updateRole(memberId: string, role: string): Promise<TeamMember> {
    const response = await this.apiClient.patch<TeamMember>(
      `/team/members/${memberId}`,
      { role }
    );
    return response.data;
  }
}
